var http = require('http');

http.createServer(function(request, response) {
   console.log('Received request for URL: ' + request.url);
   response.writeHead(200);
   response.end('Hello World');
}).listen(8080);
