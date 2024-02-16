package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
)

func SendRequest(port string, pattern string, sendData interface{}, responseDataArr...interface{}) {
	sendDataCoded := sendData
	if _, isReadyToSend := sendData.(io.Reader); !isReadyToSend {
		jsonCoded, err := json.Marshal(sendData)
		if err != nil {
			panic("Ошибка парса в JSON")
		}
		sendDataCoded = bytes.NewReader(jsonCoded)
	}
	
	response, err := http.Post(fmt.Sprintf("http://localhost:%s/%s", port, pattern), "application/json", sendDataCoded.(io.Reader))
	if err != nil {
		panic("Ошибка запроса!")
	}
	defer response.Body.Close()

	if len(responseDataArr) != 0 {
		responseData := responseDataArr[0]
		switch x := responseData.(type) {
		case *[]byte:
			*x, err = io.ReadAll(response.Body)
			if err != nil {
				panic("Ошибка при чтении ответа без распарса! []byte")
			}
		default:
			err = json.NewDecoder(response.Body).Decode(responseData)
			if err != nil {
				panic("Ошибка распарса!")
			}
		}
		
	}
}
