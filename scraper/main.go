package main

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"google.golang.org/genai"
)

var variant string = "gemini-2.5-flash-preview-05-20"
var systemPrompt string = ` collect this info from url if available. return as csv row;
Company Name
Website URL
Industry / Sector
Description / Services
Contact Info
Social Enterprise Status
Related News/Updates
Program Participation (if available)
`

var API_COUNT int

func main() {
	apiKey := os.Getenv("API")
	client, err := genai.NewClient(context.TODO(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	// prop_read_file := map[string]*genai.Schema{
	// 	"path": &genai.Schema{
	// 		Type: genai.TypeString,
	// 	},
	// }
	// prop_list_files := map[string]*genai.Schema{
	// 	"path": &genai.Schema{
	// 		Type: genai.TypeString,
	// 	},
	// }
	// prop_edit_file := map[string]*genai.Schema{
	// 	"path": &genai.Schema{
	// 		Type: genai.TypeString,
	// 	},
	// 	"content": &genai.Schema{
	// 		Type: genai.TypeString,
	// 	},
	//}

	urlContext := &genai.URLContext{}

	tools := []*genai.Tool{
		//{
		// 	FunctionDeclarations: []*genai.FunctionDeclaration{
		// 		{
		// 			Name: "read_file",
		// 			Parameters: &genai.Schema{
		// 				Type:        genai.TypeObject,
		// 				Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
		// 				Properties:  prop_read_file,
		// 			},
		// 		},
		// 		{
		// 			Name: "list_files",
		// 			Parameters: &genai.Schema{
		// 				Type:        genai.TypeObject,
		// 				Description: "List the contents of a given relative directory path. Use this when you want to see what's inside a directory.",
		// 				Properties:  prop_list_files,
		// 			},
		// 		},
		// 		{
		// 			Name: "edit_file",
		// 			Parameters: &genai.Schema{
		// 				Type:        genai.TypeObject,
		// 				Description: "Edit the contents of a given relative file path. Use this when you want to modify a file.",
		// 				Properties:  prop_edit_file,
		// 			},
		// 		},
		// 	},
		// },
		{
			URLContext: urlContext,
		},
	}
	agent := NewAgent(client, getUserMessage, tools)
	err = agent.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func NewAgent(client *genai.Client, getUserMessage func() (string, bool), tools []*genai.Tool) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
	}
}

type Agent struct {
	client         *genai.Client
	getUserMessage func() (string, bool)
	tools          []*genai.Tool
}

func (a *Agent) Run(ctx context.Context) error {
	conversation := []*genai.Content{}
	fmt.Println("Chat with Gemini (use 'ctrl-c' to quit)")

	readUserInput := true
	for {
		if readUserInput {
			fmt.Print("\u001b[94mYou\u001b[0m: ")
			userInput, ok := a.getUserMessage()
			if !ok {
				break
			}

			userMessage := genai.NewContentFromText(userInput, "user")
			conversation = append(conversation, userMessage)
		}

		message, err := a.runInference(ctx, conversation)
		if err != nil {
			return err
		}
		conversation = append(conversation, message)
		toolConv := []*genai.Content{}
		for _, part := range message.Parts {
			if part.FunctionCall != nil {
				fmt.Printf("\u001b[93mFunCall\u001b[0m: %s\n", part.FunctionCall.Name)
				for k, v := range part.FunctionCall.Args {
					fmt.Printf("\u001b[93mKey\u001b[0m: %s\n", k)
					fmt.Printf("\u001b[93mArg\u001b[0m: %s\n", v)
				}
				data := a.executeTool(part.FunctionCall)
				toolMessage := genai.NewContentFromText(data, "user")
				toolConv = append(toolConv, toolMessage)
			} else {
				// TODO: IMPORTANT! append to existing data.csv
				edit_file("test.csv", part.Text)
				fmt.Printf("\u001b[93mGemini\u001b[0m: %s\n", part.Text)
				break
			}
		}
		if len(toolConv) == 0 {
			readUserInput = true
		} else {
			log.Printf("Tool len: %v", len(toolConv))
			conversation = append(conversation, toolConv...)
			toolConv = toolConv[:0]
			readUserInput = false
		}
	}
	return nil
}

func (a *Agent) runInference(ctx context.Context, conversation []*genai.Content) (*genai.Content, error) {
	API_COUNT++
	log.Printf("Calling API count %v", API_COUNT)
	result, err := a.client.Models.GenerateContent(ctx, variant, conversation, &genai.GenerateContentConfig{
		Tools:             a.tools,
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: systemPrompt}}},
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		log.Println("result is nil")
	}
	return result.Candidates[0].Content, err
}

func (a *Agent) executeTool(tool *genai.FunctionCall) string {
	switch tool.Name {
	case "read_file":
		data, err := read_file(tool.Args["path"].(string))
		if err != nil {
			return fmt.Sprintf("Error reading file: %v", err)
		}
		return data
	case "list_files":
		loc := tool.Args["path"]
		locPath := "."
		if loc != nil {
			locPath = loc.(string)
		}
		files, err := list_files(locPath)
		if err != nil {
			return fmt.Sprintf("Error listing files: %v", err)
		}
		return strings.Join(files, "\n")
	case "edit_file":
		path := tool.Args["path"].(string)
		content := tool.Args["content"].(string)
		err := edit_file(path, content)
		if err != nil {
			return fmt.Sprintf("Error editing file: %v", err)
		}
		return "OK"
	default:
		return fmt.Sprintf("Unknown tool: %s", tool.Name)
	}
}

func read_file(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func list_files(dirPath string) ([]string, error) {
	originalDirPath := dirPath
	if !filepath.IsAbs(dirPath) && filepath.Base(dirPath) == dirPath && dirPath != "." && dirPath != ".." {
		dirPath = filepath.Join(".", dirPath)
	}

	var files []string

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) && path == dirPath {
				return fmt.Errorf("failed to read directory '%s': %w", originalDirPath, err)
			}
			// Propagate other errors to stop the walk
			return err
		}

		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			relativePath, relErr := filepath.Rel(dirPath, path)
			if relErr != nil {
				return fmt.Errorf("could not get relative path for %q from %q: %w", path, dirPath, relErr)
			}

			if relativePath != "." {
				files = append(files, relativePath)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func edit_file(path string, newContent string) error {
	log.Println("Editing file:", path)
	_, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) && path != "" {
			err = createNewFile(path, newContent)
			if err != nil {
				return err
			}
		}
		return err
	}
	err = os.WriteFile(path, []byte(newContent), 0644)
	if err != nil {
		log.Println(err)
	}
	return err
}

func createNewFile(filePath, content string) error {
	log.Println("Creating new file:", filePath)
	dir := path.Dir(filePath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}
