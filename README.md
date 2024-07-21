### What does this fork add ?
It enables the use of goroutines to retrieve streams concurrently by mitigating concurrent writes errors.
This particular branch offers slightly better performances in addition to concurrency, as it is making less write operation (to a map we deleted).

Changes made in file xtream-codes.go :

- Remove streams map variable from XtreamClient struct :
  
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

      //FORK : we delete streams map.
  }
  ```
- Modified NewClient function :

    ```go
  func NewClient(username, password, baseURL string) (*XtreamClient, error) {

      _, parseURLErr := url.Parse(baseURL)
      if parseURLErr != nil {
        return nil, fmt.Errorf("error parsing url: %s", parseURLErr.Error())
      }

      client := &XtreamClient{
        Username:  username,
        Password:  password,
        BaseURL:   baseURL,
        UserAgent: defaultUserAgent,

        HTTP:    http.DefaultClient,
        Context: context.Background(),

        //FORK no need to initialize streams map
      }

      authData, authErr := client.sendRequest("", nil)
      if authErr != nil {
        return nil, fmt.Errorf("error sending authentication request: %s", authErr.Error())
      }

      a := &AuthenticationResponse{}

      if jsonErr := json.Unmarshal(authData, &a); jsonErr != nil {
        return nil, fmt.Errorf("error unmarshaling json: %s", jsonErr.Error())
      }

      client.ServerInfo = a.ServerInfo
      client.UserInfo = a.UserInfo

      return client, nil
  }
    ```

- Modified GetStreamURL function :
  
  ```go
  func (c *XtreamClient) GetStreamURL(stream_elem interface{}, wantedFormat string) (string, error) {
      validFormat := false

      for _, allowedFormat := range c.UserInfo.AllowedOutputFormats {
        if wantedFormat == allowedFormat {
          validFormat = true
        }
      }

      if !validFormat {
        return "", fmt.Errorf("%s is not an allowed output format", wantedFormat)
      }

      //FORK-ADDED : type assertion of stream_elem to return the correct url.
      switch stream := stream_elem.(type) {
      case Stream:
        if stream.Type == "movie" {
          return fmt.Sprintf("%s/%s/%s/%s/%d.%s", c.BaseURL, stream.Type, c.Username, c.Password, stream.ID, stream.ContainerExtension), nil
        } else {
          return fmt.Sprintf("%s/%s/%s/%s/%d.%s", c.BaseURL, stream.Type, c.Username, c.Password, stream.ID, wantedFormat), nil
        }

      case SeriesEpisode:
        return fmt.Sprintf("%s/series/%s/%s/%s.%s", c.BaseURL, c.Username, c.Password, stream.ID, stream.ContainerExtension), nil
      }

      return "", fmt.Errorf("could not determine stream url")

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

      //FORK : No need to log streams into XtreamClient struct anymore.

      return streams, nil
  }
  ```

# go.xtream-codes
A Go library for accessing Xtream-Codes servers
