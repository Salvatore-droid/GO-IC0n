package main

import (
	"fmt"
	"strings"
	"os"
	"os/exec"
	"bufio"
	"bytes"
	"text/template"
	"text/tabwriter"
	"github.com/fatih/color"
	"github.com/common-nighthawk/go-figure"
)

func main(){
	yellow:=color.New(color.FgYellow).SprintFunc()
	green:=color.New(color.FgGreen).SprintFunc()
	banner()
	asterisk()
	fmt.Println(green(`
			[[=>>Usage<<=]]

	[=>>List for flatpak downloaded apps<<=]

	WPS Office      com.wps.Office  11.1.0.11719    stable  system
	Bible   net.lugsole.bible_gui   0.1.4   stable  system
	Calculator      org.gnome.Calculator    46.1    stable  system

	[=set icon name: (example) wps=]
	[=Enter appID of app: [see (Application ID) from the list in your console] (example) com.wps.Office=]
	[=Enter path to image: [insert path to where image of icon is located] (example) /usr/=]
	`))
	asterisk()
	list_exit()
	create_icon()
	fmt.Println(yellow("Icon created successfully...!!"))
	
}


func list_exit(){
	red:=color.New(color.FgRed).SprintFunc()
	cyan:=color.New(color.FgCyan).SprintFunc()
	fmt.Println(cyan("\n[0] To exit"))
	fmt.Println(cyan("[1] To list downloaded flatpak apps\n"))

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter value: ")
	value, err :=reader.ReadString('\n')
	if err != nil{
		fmt.Println("Error:", err)
	}
	value = strings.TrimSuffix(value, "\n")

	if value =="0"{
		os.Exit(1)
	}else if value == "1"{
		app_list()
	}else{
		fmt.Println(red("Invalid input, enter 0 or 1"))
		os.Exit(0)
	}
}

func banner(){
	cyan:=color.New(color.FgCyan).SprintFunc()
	yellow:=color.New(color.FgYellow).SprintFunc()
	figure:= figure.NewFigure("PY-IC0n","", true)
	fmt.Println(cyan(figure))
	fmt.Println(yellow("\n[=[Icon maker for flatpak apps]=]"))
	fmt.Println(yellow("[=[Created by Genius]=]"))
	fmt.Println(yellow("[=[https://github.com/Salvatore-droid]=]\n"))
}


func asterisk(){
	blue:=color.New(color.FgBlue).SprintFunc()
	line:=strings.Repeat("*", 100)
	fmt.Println(blue(line))
}


func app_list(){
	green:=color.New(color.FgGreen).SprintFunc()
	blue:=color.New(color.FgBlue).SprintFunc()
	red:=color.New(color.FgRed).SprintFunc()
	yellow:=color.New(color.FgYellow).SprintFunc()
	fmt.Println(blue("\n[ Enter quit to exit ]\n"))
	fmt.Println(green("\n[=>>List of flatpak downloaded apps<<=]\n"))
	list := exec.Command("flatpak", "list", "--app")

	output, err := list.CombinedOutput()

	if err != nil{
		fmt.Println("Error:", err)
	}
	writer:=tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ',tabwriter.Debug)
	fmt.Fprintln(writer, red("Name\tApplication ID\tVersion\tBranch\tInstallation"))

	lines := strings.Split(string(output), "\n")

	for _, line := range lines{
		if line != ""{
			fmt.Fprintln(writer, yellow(line))
		}
	}

	writer.Flush()

}
func create_icon(){
	red:=color.New(color.FgRed).SprintFunc()
	const CommandTemplate =`echo "[Desktop Entry]
		Version=1.0
		Name={{.Name}}
		Comment=App
		Type=Application
		Exec=flatpak run {{.AppID}}
		Icon={{.Path}}
		Terminal=false
		Categories=Utility;" > ~/Desktop/{{.Name}}.desktop
		`

	type UserInput struct{
		Name string
		AppID string
		Path string
	}

	reader :=bufio.NewReader(os.Stdin)

	fmt.Print("\nset name of icon: ")
	name, err :=reader.ReadString('\n')
	if err != nil{
		fmt.Println("Error reading input:", err)			
		os.Exit(1)
	}

	name=strings.TrimSuffix(name, "\n")

	if name == "quit"{
		fmt.Println(red("exiting program..."))
		os.Exit(0)
	}
	
	fmt.Print("Enter appID of app: ")
	appID, err := reader.ReadString('\n')
	if err != nil{
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	appID = strings.TrimSuffix(appID, "\n")

	if appID == "quit"{
		fmt.Println(red("exiting program..."))
		os.Exit(0)
	}


	fmt.Print("Enter the path to image(optional, press enter to skip):")
	path, err :=reader.ReadString('\n')
	if err != nil{
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	path = strings.TrimSuffix(path, "\n")

	input := UserInput{Name: name, AppID: appID, Path: path}

	t := template.Must(template.New("cmd").Parse(CommandTemplate))
	var cmdBuffer bytes.Buffer
	t.Execute(&cmdBuffer, input)

	cmdStr := cmdBuffer.String()

	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil{
		fmt.Println("Error:", err)
	}else{
		fmt.Println(string(output))
	}

}

