package controllers

//ApplicationData describes the connection information about an application you can run on a chromecast.
type ApplicationData struct {
	ID        string //ID used for application to be launched/manipulated on chromecast.
	Namespace string //Namespace for communication channel for application on chromecast.
}
