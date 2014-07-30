package main

import (
	"fmt"

	"github.com/mitchellh/colorstring"
)

func main() {
	say("[blue]> Hold still citizen, scanning dependencies for contraband...")
	say("[white]github.com/xoebus/apache                                    [green]CHECKS OUT")
	say("[white]github.com/xoebus/copyright                                 [red]CONTRABAND")
	say("[white]github.com/xoebus/mit                                       [green]CHECKS OUT")
	say("")

	say("[blue]> We found questionable material. Citizen, what do you have to say for yourself?")
	say("[white]github.com/xoebus/no-license                                [magenta]NO LICENSE")
	say("[white]github.com/xoebus/greylist                                  [yellow]BORDERLINE")
}

func say(message string) {
	fmt.Println(colorstring.Color(message))
}
