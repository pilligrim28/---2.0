import express from 'express';
import { Low } from 'lowdb';
import { JSONFile } from 'lowdb/node';
import { createServer } from 'http';
import { Server } from 'socket.io';
import dgram from 'dgram';
import net from 'net';
import sip from 'sip';

// JSON-база
const adapter = new JSONFile('db.json');
const db = new Low(adapter);
await db.read();
db.data ||= { radios: [] };

// Express
const app = express();
const httpServer = createServer(app);
const io = new Server(httpServer);

// TCP-сервер
const tcpServer = net.createServer(socket => {
  socket.on('data', data => {
    console.log('TCP:', data.toString());
    // Логика обработки данных от рации
  });
});

// UDP-сервер
const udpServer = dgram.createSocket('udp4');
udpServer.on('message', (msg, rinfo) => {
  console.log('UDP:', msg.toString());
});

// SIP обработчик
sip.start({}, (req, res) => {
  console.log('SIP:', req);
});

// API
app.use(express.json());
app.get('/radios', (req, res) => res.json(db.data.radios));
app.post('/radios', async (req, res) => {
  db.data.radios.push(req.body);
  await db.write();
  io.emit('radio-added', req.body); // Вещание через WebSocket
  res.sendStatus(201);
});

// Запуск
tcpServer.listen(5000, '0.0.0.0');
udpServer.bind(5001);
httpServer.listen(3000);
console.log('Сервер запущен на портах 3000 (HTTP), 5000 (TCP), 5001 (UDP)');