package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

func main() {

	inDirPtr := flag.String("i", "./", "Input Directory or File (*.smali)")
	outDirPtr := flag.String("o", "java/", "Output Directory")
	flag.Parse()

	err := filepath.Walk(*inDirPtr,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(*inDirPtr, path)
			if err != nil {
				return err
			}

			if info.IsDir() {
				os.Mkdir(filepath.Join(*outDirPtr, relPath), info.Mode().Perm())
			} else if filepath.Ext(path) == ".smali" {
				maxLine, err := maxLineFromSmali(path)
				if err != nil {
					return err
				}

				//log.Println(path, maxLine)

				err = writeBogusJava(maxLine, *outDirPtr, relPath)
				if err != nil {
					return err
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func maxLineFromSmali(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	maxLine := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		fields := strings.Fields(line)
		if len(fields) > 1 && fields[0] == ".line" {
			lineNum, err := strconv.Atoi(fields[1])
			if err != nil {
				return 0, err
			}
			if lineNum > maxLine {
				maxLine = lineNum
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return maxLine, nil
}

type Tmpl struct {
	Package  string
	Filename string
}

func writeBogusJava(maxLine int, basedir string, relpath string) error {
	file, err := os.Create(filepath.Join(basedir, relpath))
	if err != nil {
		return err
	}
	defer file.Close()

	header, err := template.New("header").Parse(`
package {{.Package}};

public class {{.Filename}} {
  public int Function() {
    int x = 0;`)
	if err != nil {
		return err
	}
	filename := filepath.Base(relpath)
	err = header.Execute(file,
		Tmpl{Package: strings.ReplaceAll(filepath.Dir(relpath), string(filepath.Separator), "."),
			Filename: strings.TrimSuffix(filename, filepath.Ext(filename))})
	if err != nil {
		return err
	}
	file.WriteString("\n")
	for i := 0; i < maxLine; i++ {
		file.WriteString("    x += 1; //" + strconv.Itoa(i) + "\n")
	}
	file.WriteString("    return x;\n  }\n};\n")

	return nil
}
