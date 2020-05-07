package main

import "github.com/vlad-m-r/checker/api/controllers"

func main() {
	controllers.RunChecks("config.yaml")
}
