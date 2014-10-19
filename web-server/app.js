var express    = require('express');
var request = require("request")
var bodyParser = require('body-parser');

var app = express();

app.set('view engine', 'jade')
app.set('views', __dirname + '/views')
var router = express.Router();

app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());

var port = 14000;
var host = "";
var containers = [];




function getContainers() {
	request({
		url: "unix:///var/run/docker.sock/containers/json",
		json: true
		}, function (err, res, body) {
		if (!err && res.statusCode == 200) {
			containers = body
		} else if (err) {
			console.log(err);
		}
	});
}

function createContainer(createPort, joinPort) {
	var addr = host+":"+createPort;
	var joinaddr = host+":"+joinPort;
	var tcp = createPort + "/tcp";
	var udp = createPort + "/udp";
	
	var exportedPorts = {};
    exportedPorts[tcp] = {};
	exportedPorts[udp] = {};
	
	var json = {
		"Cmd": ["-l",
			addr,
			"-j",
			joinaddr],
		"Image":"d7024e",
		"ExposedPorts": exportedPorts
};
	
	request.post({
		headers: { 'content-type': 'application/json' },
		url: "unix:///var/run/docker.sock/containers/create",
		json: json
	}, function(err, res, body){
		if (!err) {
			console.log(body)
			startContainer(body.Id, createPort);
		}
	});
}

function startContainer(id, port) {
	var tcp = port + "/tcp"
	var udp = port + "/udp"
	
	var portBindings = {};
    portBindings[tcp] = [{ "HostPort": port }];
	portBindings[udp] = [{ "HostPort": port }];
	
	var json = {
		"LxcConf": null,
		"NetworkMode": "bridge",
		"PortBindings": portBindings
	};
	request.post({
		headers: { 'content-type': 'application/json' },
		url: "unix:///var/run/docker.sock/containers/"+id+"/start",
		json: json
	}, function(error, res, body){
		console.log(body)
	});
}

function stopContainer(id) {
	request.post({
		url: "unix:///var/run/docker.sock/containers/"+id+"/stop"
	}, function(err, res, body){
		console.log(body)
	});
}

function deleteContainer(id) {
	request({
		method: 'DELETE',
		url: "unix:///var/run/docker.sock/containers/"+id
	}, function(err, res, body){
		console.log(body)
	});
}

app.get('/', function (req, res) {
  res.render('index',
  { title : 'Home' }
  )
})

router.get('/containers', function(req, res) {
	// GET all containers and oouput them
	getContainers();
	
	var output = "<p><a href='/'>Back</a><p>"
	
	// <pre> makes the text appear as plain text instead of html
	output += "<pre>#   ID                                                                 PORT\n"
	
	for (var i = 0; i < containers.length; i++) {
		output += i + "   " + containers[i]["Id"] + "   " + containers[i]["Ports"][0].PublicPort + "\n";
	}
	res.send(output);
});

app.post('/create', function(req, res) {
	var createPort = req.body.createPort;
	var joinPort = req.body.joinPort;
	createContainer(createPort, joinPort);
});

app.post('/stop', function(req, res) {
	var id = req.body.stopId;
	stopContainer(id);
});

app.post('/delete', function(req, res) {
	var id = req.body.deleteId;
	deleteContainer(id);
});

app.use('/', router);

app.listen(port);