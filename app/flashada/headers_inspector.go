package flashada

type KeyValueAPI[T any] interface {
	Get(string) (*T, error)
	Set(T) error
}

type httpHeadersAPI struct {
	headers http.Header
}

func NewHttpHeadersAPI(headers http.Header) *httpHeadersAPI {
	return &httpHeadersAPI{headers: headers}
}

func (api *httpHeadersAPI) Get(key string) (*string, error) {
	value := api.headers.Get(key)
	if value == "" {
		return nil, fmt.Errorf("header not found")
	}
	return &value, nil
}

func (api *httpHeadersAPI) Set(key string, value string) error {
	api.headers.Set(key, value)
	return nil
}type grpcMetadataAPI struct {
	md metadata.MD
}

func NewGrpcMetadataAPI(md metadata.MD) *grpcMetadataAPI {
	return &grpcMetadataAPI{md: md}
}

func (api *grpcMetadataAPI) Get(key string) (*string, error) {
	values := api.md.Get(key)
	if len(values) == 0 {
		return nil, fmt.Errorf("metadata not found")
	}
	return &values[0], nil
}

func (api *grpcMetadataAPI) Set(key string, value string) error {
	// This operation is conceptually invalid for incoming metadata in a server-side context,
	// as incoming metadata is immutable. For demonstration purposes only:
	api.md.Set(key, value)
	return nil
}type redisAPI struct {
	client *redis.Client
}

func NewRedisAPI(client *redis.Client) *redisAPI {
	return &redisAPI{client: client}
}

func (api *redisAPI) Get(key string) (*string, error) {
	val, err := api.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key does not exist")
	} else if err != nil {
		return nil, err
	}
	return &val, nil
}

func (api *redisAPI) Set(key string, value string) error {
	err := api.client.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
