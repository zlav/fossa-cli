package test

import (
	"testing"

	"github.com/fossas/fossa-cli/config"
	gock "gopkg.in/h2non/gock.v1"
)

func TestPublishWrongResponseStatus(t *testing.T) {
	defer gock.Off()

	url := "http://server.com"
	gock.New(url).
		Get("/bar").
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	config.BackendEndpoint = url
	// test.Run(&cli.Context{App:,
	// flagSet: ,
	// parentContext: ,})

	// stdout, stderr, err := exec.Run(exec.Cmd{
	// 	Name: "fossa",
	// 	Argv: append([]string{"test", "-c="}),
	// })
	// fmt.Printf("%+v\n", stdout)
	// fmt.Printf("%+v", stderr)
	// if err != nil && stdout == "" {
	// 	if strings.Contains(stderr, "build constraints exclude all Go files") {
	// 		fmt.Println(errors.New("bad OS/architecture target"))
	// 	}
	// 	fmt.Println("error running fossa test")
	// }
}
