### What does this fork add ?
It enables the use of goroutines to retrieve streams concurrently by mitigating concurrent writes errors.

Changes made in file xtream-codes.go :

- New import :
  
  ```go
  import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strconv"
    "sync" // <- New import 
  )
  ```
- New variable in the XtreamClient struct :

    ```go
  type XtreamClient struct {
      Username  string
      Password  string
      BaseURL   string
      UserAgent string
    
      ServerInfo ServerInfo
      UserInfo   UserInfo
    
      // Our HTTP client to communicate with Xtream
      HTTP    *http.Client
      Context context.Context
    
      // We store an internal map of Streams for use with GetStreamURL
      streams map[int]Stream
    
      //FORK-ADDED : We add a mutex to mitigate concurrency writes errors to the streams map.
      mu sync.Mutex // <- New mutex variable
  }
    ```

- Modified GetStreams function :
  
  ```go
  func (c *XtreamClient) GetStreams(streamAction, categoryID string) ([]Stream, error) {
	var params url.Values
	if categoryID != "" {
		params = url.Values{}
		params.Add("category_id", categoryID)
	}

	// For whatever reason, unlike live and vod, series streams action doesn't have "_streams".
	if streamAction != "series" {
		streamAction = fmt.Sprintf("%s_streams", streamAction)
	}

	streamData, streamErr := c.sendRequest(fmt.Sprintf("get_%s", streamAction), params)
	if streamErr != nil {
		return nil, streamErr
	}

	streams := make([]Stream, 0)

	if jsonErr := json.Unmarshal(streamData, &streams); jsonErr != nil {
		return nil, jsonErr
	}

	//FORK-ADDED : Locking access to streams map while new values are insterted .
	//It is now secure to use multiple goroutines to fetch streams.
	c.mu.Lock() // <- added line
	for _, stream := range streams {
		c.streams[int(stream.ID)] = stream
	}
	c.mu.Unlock() // <- added line
	//Unlocking access to streams map.

	return streams, nil
  }
  ```

# go.xtream-codes
A Go library for accessing Xtream-Codes servers
