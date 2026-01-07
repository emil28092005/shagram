let ws = null;
const messagesDiv = document.getElementById('messages');
const messageInput = document.getElementById('messageInput');
const roomInput = document.getElementById('roomInput');

function connectRoom() {
    const room = roomInput.value || 'general';
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//localhost:8080/ws/${room}`);

    ws.onopen = function() {
        console.log('Connected to room: ' + room);
        messagesDiv.innerHTML = '<div class="message">✅ Connected to ' + room + '</div>';
    };

    ws.onmessage = function(event) {
        const msgDiv = document.createElement('div');
        msgDiv.className = 'message';
        msgDiv.textContent = event.data;
        messagesDiv.appendChild(msgDiv);
        messagesDiv.scrollTop = messagesDiv.scrollHeight;
    };

    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
        messagesDiv.innerHTML += '<div class="message">❌ Error: ' + error + '</div>';
    };

    ws.onclose = function() {
        messagesDiv.innerHTML += '<div class="meassage">⛔ Disconnected</div>';
    };
}

function sendMessage() {
    const text = messageInput.value.trim();
    if (!text || !ws || ws.readyState !== WebSocket.OPEN) {
        alert('Not connected or empty message');
        return;
    }
    
    ws.send(JSON.stringify({text: text}));
    messageInput.value = '';
}