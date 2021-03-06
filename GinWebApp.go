package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	Url = "http://localhost:8181/checkText"
	BodyType = "application/json"
	BadResponse = "HTTP Code 204 No Content"
	Port = ":8181"
)

type Request struct {
	Site []string
	SearchText string
}
func (r *Request )Init(_Site []string, _SearchText string){
	r.Site = _Site
	r.SearchText = _SearchText
}

type Response struct {
	FoundAtSite string
}
func CheckError(msg string, e *error ){
	if *e != nil{
		fmt.Println(msg,*e)
	}
}
//second parametr ONLY pointer
func DecodeJSON(from [] byte, to interface{})error{
	return json.NewDecoder(bytes.NewBuffer(from)).Decode(to)
}

func server(){
	gin.SetMode(gin.ReleaseMode)
	router:= gin.Default()
	router.POST("/checkText",func (c *gin.Context){
		request :=new(Request)
		response := new(Response)
		requestBodyJSON,_:= ioutil.ReadAll(c.Request.Body)
		err := DecodeJSON(requestBodyJSON,request)
		CheckError("json was not decode",&err)
		response.FoundAtSite = BadResponse
		for _, str := range request.Site {
			if strings.Contains(strings.ToLower(str), strings.ToLower(request.SearchText)) {
				response.FoundAtSite = str
				break
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"FoundAtSite": response.FoundAtSite,
		})
	})
	fmt.Println("Server started on port " + Port)
	router.Run(Port)
}


func client(request *Request){
	response:= new(Response)
	requestBodyJSON, err := json.Marshal(request)
	CheckError("json was not encode",&err)
	responseBytes,err:=http.Post(Url, BodyType, bytes.NewBuffer(requestBodyJSON))
	CheckError("bad request or no connection",&err)
	responseDataJSON, err := ioutil.ReadAll(responseBytes.Body)
	CheckError("No response",&err)
	err = DecodeJSON(responseDataJSON,response)
	CheckError("json was not decode",&err)
	fmt.Println("Response: ",response.FoundAtSite)
}
func main() {
	go server()
	request:= new(Request)
	request.Init([]string {"https://google.com", "https://yahoo.com"},"lol")
	client(request)
	request.SearchText ="yahoo"
	client(request)
	request.SearchText = "Google"
	client(request)
}
