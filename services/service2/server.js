const http = require('http');

const server = http.createServer((req, res) => {
    if (req.url === '/resources') {
        res.writeHead(200, { 'Content-Type': 'text/plain' });
        res.end('Hello from Service 2!');
    }
});

server.listen(8082, () => {
    console.log('Service 2 running on port 8082');
});