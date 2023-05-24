package store

func DoStorePut(key string, user string, value any) StorePutResponse {
	responseChannel := make(chan StorePutResponse)
	request := StorePutRequest{
		Key:         key,
		User:        user,
		Data:        value,
		RespChannel: responseChannel,
	}

	requestChannel <- request
	return <-responseChannel
}

func DoStoreGet(key string) StoreGetResponse {
	responseChannel := make(chan StoreGetResponse)
	request := StoreGetRequest{
		Key:         key,
		RespChannel: responseChannel,
	}

	requestChannel <- request
	return <-responseChannel
}
