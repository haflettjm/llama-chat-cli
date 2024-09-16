package main

func main() {
	networkChan := make(chan Result)
	go func() {
		resp := sendChat()
		result := parseResp(resp)
		resultChan <- result

	}()
	process(<-networkChan, <-networkChan)
}
