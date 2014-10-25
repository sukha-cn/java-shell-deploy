package main;

import (
	"os";
	"os/exec";
	"fmt";
	"net/http";
	"io/ioutil";
	"encoding/json";
	"strings";
	"strconv";
	"flag";
);

type AppContext struct {
	UrlPattern string;
	Port string;
	SecretKey string;
	SecretValue string;
	ScriptPath string;
	ScriptParameter string;
	LogFile string;
};

var context AppContext;

func castToBool(value interface{}) bool {

	boolValue, ok := value.(bool);

	if(!ok) {
		return false;
	}

	return boolValue;

}

func castToInt(value interface{}) int {

	intValue, ok := value.(int);

	if(!ok) {
		return 0;
	}

	return intValue;

}

func castToFloat(value interface{}) float64 {

	floatValue, ok := value.(float64);

	if(!ok) {
		return 0;
	}

	return floatValue;

}

func castToString(value interface{}) string {

	stringValue, ok := value.(string);

	if(!ok) {
		return "";
	}

	return stringValue;

}

func ensureJsonValue(jsonInterface interface{}, path string, expectedValue string) bool {

	var splittedPath = strings.Split(path, "/");

	var jsonElement = jsonInterface.(map[string]interface{});
	var result = "";

	for _,element := range splittedPath {
		value,ok := jsonElement[element];
		if(!ok) {
			break;
		}

		switch t := value.(type) {
			case bool:
			    castedValue := castToBool(value);
			    result = strconv.FormatBool(castedValue);
			    break;
			case int:
			    castedValue := castToInt(value);
			    result = strconv.Itoa(castedValue);
			    break;
			case float64:
			    castedValue := castToFloat(value);
			    result = strconv.FormatFloat(castedValue, 'f', -1, 64);
			    break;
			case string:
			    castedValue := castToString(value);
			    result = castedValue;
			    break;
			case interface{}:
			    jsonElement = value.(map[string]interface{});
			default:
			    fmt.Printf("unexpected type %T", t)
		}				
	}

	return result == expectedValue;

}

func logToFile(message string) {

	fmt.Println(fmt.Sprintf("%s: %s", context.LogFile, message));

	file, err := os.OpenFile(context.LogFile, os.O_APPEND|os.O_WRONLY, 0600);
	if (err != nil) {
		file, err = os.Create(context.LogFile);
		if (err != nil) {
	        panic(err)
	    }
    }

    defer file.Close();

    _, err = file.WriteString(message); 

    if (err != nil) {
    	panic(err)
    }

}

func executeShellScript() string {

    var cmd = exec.Command(context.ScriptPath, context.ScriptParameter);
    out, err := cmd.Output();

    if (err != nil) {
        logToFile(err.Error());
        return err.Error();
    }

    var result = string(out);
    logToFile(result);
    return string(result);

}

func rootHandler(response http.ResponseWriter, request *http.Request) {

	fmt.Fprintf(response, request.URL.Path);

}

func notificationHandler(response http.ResponseWriter, request *http.Request) {

	logToFile("------------------\n");
	logToFile(fmt.Sprintf("Content-Type: %s, Event: %s, Delivery: %s\n", 
		request.Header["Content-Type"], request.Header["X-Github-Event"], request.Header["X-Github-Delivery"]));

    bodyBytes, err := ioutil.ReadAll(request.Body);
    if (err != nil) {
        panic(err);
    }

    var jsonInterface interface{};
    var result = "";
    
	jsonErr := json.Unmarshal(bodyBytes, &jsonInterface);
    if (jsonErr != nil) {
        result = jsonErr.Error();
    } else {
	    if(ensureJsonValue(jsonInterface, context.SecretKey, context.SecretValue)) {
	    	logToFile(executeShellScript());
	    	result = "done\n";
    	} else {
    		result = "not the right hook\n";
    	}
    }

    logToFile(result);
    fmt.Fprintf(response, result);

	logToFile(fmt.Sprintf("Body: %s\n", string(bodyBytes))); 
	logToFile("------------------\n");

}


func main() {

	var pattern = flag.String("pattern", "test", "a string");
	var port = flag.String("port", "8081", "a string");
	var secretKey = flag.String("key", "key", "a string");
	var secretValue = flag.String("value", "secret", "a string");
	var scriptPath = flag.String("script", "service.sh", "a string");
	var scriptParameter = flag.String("parameter", "redeploy", "a string");
	var logFile = flag.String("log", "log.txt", "a string");

	flag.Parse();

	context = AppContext{
		UrlPattern: *pattern,
		Port: *port,
		SecretKey: *secretKey,
		SecretValue: *secretValue,
		ScriptPath: *scriptPath,
		ScriptParameter: *scriptParameter,
		LogFile: *logFile};

	http.HandleFunc("/", rootHandler);
	http.HandleFunc("/" + context.UrlPattern, notificationHandler);
	http.ListenAndServe(":" + context.Port, nil);

}