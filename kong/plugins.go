package kong

import (
	"fmt"
	"github.com/fatih/structs"
	"net/http"
	"reflect"
	"strings"
)

// PluginsService handles communication with Kong's '/plugins' resource
type PluginsService service

// Plugins represents the object returned from Kong when querying for
// multiple plugin objects.
//
// In cases where the number of objects returned exceeds the maximum,
// Next holds the URI for the next set of results.
// i.e. "http://localhost:8001/plugins?size=2&offset=4d924084-1adb-40a5-c042-63b19db421d1"
type Plugins struct {
	Data  []Plugin `json:"consumer,omitempty"`
	Total int      `json:"total,omitempty"`
	Next  string   `json:"next,omitempty"`
}

// Plugin represents a single Kong plugin object.
//
// Because there are a myriad of structures for the
// Kong plugin config's, Plugin.Config is the generic
// map[string]interface{}. Helper method's exist on the
// more specific Plugin definitions to convert them to/from
// this struct.
type Plugin struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	CreatedAt  int                    `json:"created_at,omitempty"`
	Enabled    bool                   `json:"enabled,omitempty"`
	ApiID      string                 `json:"api_id,omitempty"`
	ConsumerID string                 `json:"consumer_id,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
}

// PluginsService.Get queries for a single Kong plugin object by id
//
// Equivalent to GET /plugins/{id}
func (s *PluginsService) Get(id string) (*Plugin, *http.Response, error) {
	u := fmt.Sprintf("plugins/%v", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(Plugin)
	resp, err := s.client.Do(req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

// EnabledPlugins represents the list of Plugins returned
// when querying /plugins/enabled
type EnabledPlugins struct {
	Plugins []string `json:"enabled_plugins,omitempty"`
}

// PluginsService.GetEnabled queries for the list of all enabled
// Kong plugins.
func (s *PluginsService) GetEnabled() (*EnabledPlugins, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "plugins/enabled", nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(EnabledPlugins)
	resp, err := s.client.Do(req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

// PluginsService.Patch updates an existing Kong plugin object for
// a specific api. Accepts either api name or id.
//
// Equivalent to PATCH /apis/{name or id}/plugins/{id}
func (s *PluginsService) Patch(api string, plugin *Plugin) (*http.Response, error) {
	u := fmt.Sprintf("apis/%v/plugins/%v", api, plugin.ID)

	req, err := s.client.NewRequest("PATCH", u, plugin)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

// PluginsService.Delete deletes a single Kong plugin object attached
// to a specifc api. Accepts either api name or id. Only accepts
// plugin id.
//
// Equivalent to DELETE /apis/{name or id}/plugins/{id}
func (s *PluginsService) Delete(api string, plugin string) (*http.Response, error) {
	u := fmt.Sprintf("apis/%v/plugins/%v", api, plugin)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// PluginsService.Post creates a new Kong plugin object.
// Which consumer and api objects the plugin gets applied to
// depend on the values of ConsumerID and ApiID on the
// passed plugin object.
//
// For more info see:
// https://getkong.org/docs/0.9.x/admin-api/#add-plugin
func (s *PluginsService) Post(plugin *Plugin) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "plugins", plugin)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

// PluginsGetAllOptions specifies optional filter parameters
// to the PluginsService.GetAll method.
//
// Additional information about filtering options can be found in
// the Kong documentation at:
// https://getkong.org/docs/0.9.x/admin-api/#list-all-plugins
type PluginsGetAllOptions struct {
	ID         string `url:"id,omitempty"`          // A filter on the list based on the id field.
	Name       string `url:"name,omitempty"`        // A filter on the list based on the name field.
	ApiID      string `url:"api_id,omitempty"`      // A filter on the list based on the api_id field.
	ConsumerID string `url:"consumer_id,omitempty"` // A filter on the list based on the consumer_id field.
	Size       int    `url:"size,omitempty"`        // A limit on the number of objects to be returned.
	Offset     string `url:"offset,omitempty"`      // A cursor used for pagination. offset is an object identifier that defines a place in the list.

}

// PluginsService.GetAll queries for all Kong plugins objects.
// This query can be filtered by supplying the PluginsGetAllOptions struct.
//
// Equivalent to GET /plugins?uri=params&from=opt
func (s *PluginsService) GetAll(opt *PluginsGetAllOptions) (*Plugins, *http.Response, error) {
	u, err := addOptions("plugins", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	plugins := new(Plugins)
	resp, err := s.client.Do(req, plugins)
	if err != nil {
		return nil, resp, err
	}

	return plugins, resp, err
}

// PluginsService.GetSchema queries for the schema of a particular
// Kong plugin.
//
// Equivalent to GET /plugins/schema/{name}
func (s *PluginsService) GetSchema(name string) (map[string]interface{}, *http.Response, error) {
	u := fmt.Sprintf("plugins/schema/%v", name)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := make(map[string]interface{})
	resp, err := s.client.Do(req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

// isZero is used when marshaling the explicit plugin types to the
// more generic Plugin.
//
// Checks if an interface{} is equal to it's underlying type's 'zero'
// value and should not be appended to Plugin.Config's map[string]interface{}
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanSet() {
				z = z && isZero(v.Field(i))
			}
		}
		return z
	case reflect.Ptr:
		return isZero(reflect.Indirect(v))
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	result := v.Interface() == z.Interface()

	return result
}

// https://getkong.org/plugins/acl/
type ACLPlugin struct {
	Plugin
	Config *ACLConfig `json:"config,omitempty"`
}

// ToPlugin converts ACLPlugin to the more generic Plugin.
// This assumes that Config does not contain nested fields
func (c *ACLPlugin) ToPlugin() *Plugin {
	s := structs.New(c.Config)
	fieldNames := s.Names()

	config := make(map[string]interface{})

	// Iterate over Config fields and create a map with keys based on the json tags
	for _, v := range fieldNames {
		tags := strings.Split(s.Field(v).Tag("json"), ",")

		// isZero checks whether the value of the field v is 'zero'
		// If it is, we do not need to add it to our map[string]interface{}
		// Doing so causes errors later when we try to marshal to JSON
		if ok := isZero(reflect.ValueOf(s.Field(v).Value())); !ok {
			config[tags[0]] = s.Field(v).Value()
		}

	}

	plugin := &c.Plugin
	plugin.Config = config
	return plugin
}

type ACLConfig struct {
	Whitelist []string `json:"whitelist,omitempty"`
	Blacklist []string `json:"blacklist,omitempty"`
}

// https://getkong.org/plugins/request-size-limiting/
type RequestSizeLimitingPlugin struct {
	Plugin
	Config RequestSizeLimitingConfig `json:"config,omitempty"`
}

type RequestSizeLimitingConfig struct {
	AllowedPayloadSize int `json:"allowed_payload_size,omitempty"`
}

// https://getkong.org/plugins/correlation-id/
type CorrelationIDPlugin struct {
	Plugin
	Config CorrelationIDConfig `json:"config,omitempty"`
}

type CorrelationIDConfig struct {
	HeaderName     string `json:"header_name,omitempty"`
	Generator      string `json:"generator,omitempty"`
	EchoDownstream bool   `json:"echo_downstream,omitempty"`
}

// https://getkong.org/plugins/rate-limiting/
type RateLimitingPlugin struct {
	Plugin
	Config RateLimitingConfig `json:"config,omitempty"`
}

type RateLimitingConfig struct {
	Second        int    `json:"second,omitempty"`
	Minute        int    `json:"minute,omitempty"`
	Hour          int    `json:"hour,omitempty"`
	Day           int    `json:"day, omitempty"`
	Month         int    `json:"month,omitempty"`
	Year          int    `json:"year, omitempty"`
	LimitBy       string `json:"limit_by,omitempty"`
	Policy        string `json:"policy,omitempty"`
	FaultTolerant bool   `json:"fault_tolerant,omitempty"`
	RedisHost     string `json:"redis_host,omitempty"`
	RedisPort     string `json:"redis_port,omitempty"`
	RedisPassword string `json:"redis_password,omitempty"`
	RedisTimeout  int    `json:"redis_timeout,omitempty"`
}

// https://getkong.org/plugins/jwt/
type JWTPlugin struct {
	Plugin
	Config JWTConfig `json:"config,omitempty"`
}

type JWTConfig struct {
	URIParamNames  []string `json:"uri_param_names,omitempty"`
	ClaimsToVerify []string `json:"claims_to_verify,omitempty"`
	KeyClaimName   string   `json:"key_claim_name,omitempty"`
	SecretIsBase64 bool     `json:"secret_is_base64,omitempty"`
}

// https://getkong.org/plugins/file-log/
type FileLogPlugin struct {
	Plugin
	Config FileLogConfig `json:"config,omitempty"`
}

type FileLogConfig struct {
	Path string `json:"path,omitempty"`
}

// https://getkong.org/plugins/key-authentication/
type KeyAuthenticationPlugin struct {
	Plugin
	Config KeyAuthenticationConfig `json:"config,omitempty"`
}

type KeyAuthenticationConfig struct {
	KeyNames        []string `json:"key_names,omitempty"`
	HideCredentials bool     `json:"hide_credentials,omitempty"`
}
