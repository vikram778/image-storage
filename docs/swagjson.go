// This file is used to generate the swagger json file
// from the api definitions stored in the definitions folder

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	SwaggerFileName       = "swagger.json"
	SwaggerHeaderFileName = "head.json"
	SwaggerComponentsName = "components.json"
)

type SwaggerHeaderInfo struct {
	Description string `json:"description"`
	Version     string `json:"version"`
	Title       string `json:"title"`
}

type SecurityDef struct {
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
	In   string `json:"in,omitempty"`
}

// SwaggerComponent ...
type SwaggerComponent struct {
	Schema map[string]map[string]interface{} `json:"schemas"`
}

type Swagger struct {
	Openapi             string                            `json:"openapi"`
	Info                SwaggerHeaderInfo                 `json:"info"`
	Host                string                            `json:"host"`
	Schemes             []string                          `json:"schemes"`
	BasePath            string                            `json:"basePath, omitempty"`
	SecurityDefinitions map[string]SecurityDef            `json:"securityDefinitions"`
	Paths               map[string]map[string]interface{} `json:"paths"`
	Components          SwaggerComponent                  `json:"components"`
}

func main() {
	var (
		defFolder, outFolder      *string
		inputFolder, outputFolder string
		d, finalDefs              []byte
		definitions               Swagger
		err                       error
		files                     []os.FileInfo
		s                         []string
	)

	defFolder = flag.String("input", "", "Path to the folder where the json file with api definitions are stored.")
	outFolder = flag.String("output", "", "Path to the folder where the final swagger.json file will be stored")
	flag.Parse()

	if *defFolder == "" {
		log.Fatal("Please provide the input folder")
	}

	inputFolder = strings.TrimRight(*defFolder, "/") + string(os.PathSeparator)

	if *outFolder == "" {
		log.Fatal("Please provide the output folder.")
	}

	outputFolder = strings.TrimRight(*outFolder, "/") + string(os.PathSeparator)

	// parse header
	if d, err = LoadFile(inputFolder + SwaggerHeaderFileName); err != nil {
		log.Fatal("Invalid header file or file missing." + err.Error())
	}

	if err = json.Unmarshal(d, &definitions); err != nil {
		log.Fatal("Invalid header file content." + err.Error())
	}

	if os.Getenv("APP_ENV") == "TESTING" {
		s = strings.Split(os.Getenv("APP_DOMAIN"), "://")
	} else {
		s = strings.Split(os.Getenv("APP_EXTERNAL_DOMAIN"), "://")
	}

	schema, host := s[0], s[1]
	definitions.Host = host
	definitions.Schemes = []string{schema}
	// get all the defs and merge them
	if files, err = ioutil.ReadDir(inputFolder); err != nil {
		log.Fatal("Error is reading input directory content.")
	}

	log.Printf("Reading api definitions from %s...\n", outputFolder)

	for _, f := range files {
		var (
			data   []byte
			apiDef Swagger
		)
		if f.IsDir() || f.Name() == SwaggerHeaderFileName {
			continue
		}

		log.Printf("Processing api definition file : %s ...\n", f.Name())

		if data, err = LoadFile(inputFolder + f.Name()); err != nil {
			log.Fatal("Invalid api definition file." + err.Error())
		}

		if err = json.Unmarshal(data, &apiDef); err != nil {
			log.Fatal("Invalid header file content." + err.Error())
		}

		if f.Name() == SwaggerComponentsName {
			definitions.Components.Schema = apiDef.Components.Schema
			continue
		}

		// check if it contains api def

		if len(apiDef.Paths) < 1 {
			log.Fatalf("Input file %s does not contain any api definitions.", f.Name())
		}

		// iterate over the api def and push to the final
		for apiPath, apiMethodDef := range apiDef.Paths {
			if _, ok := definitions.Paths[apiPath]; ok {
				// merge the api method definitions
				existingDefs := definitions.Paths[apiPath]
				for m, md := range apiMethodDef {
					existingDefs[m] = md
				}
			} else {
				definitions.Paths[apiPath] = apiMethodDef
			}
		}
	}

	// write final api def to file

	if finalDefs, err = json.Marshal(definitions); err != nil {
		log.Fatal("Unable to convert to json." + err.Error())
	}

	if err = ioutil.WriteFile(outputFolder+SwaggerFileName, finalDefs, 0766); err != nil {
		log.Fatal("Unable to write to output file." + err.Error())
	}

	log.Println("Successfully created swagger definitions.")

	os.Exit(0)

}

// LoadFile ...
func LoadFile(path string) (data []byte, err error) {

	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}

	return
}
