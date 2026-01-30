let ws = null;
let accessToken = null;

const messagesDiv = document.getElementById('messages');
const messageInput = document.getElementById('messageInput');
const roomSelect = document.getElementById('roomSelect');
const authStatus = document.getElementById('authStatus');

async function login() {
    const usernameInput = document.getElementById('userInput');
    const username = usernameInput ? usernameInput.value.trim() : '';

    if (!username) {
        alert('Enter username');
        return;
    }

    authStatus.textContent = 'Logging in...';

    const resp = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username }),
    });

    if (!resp.ok) {
        const text = await resp.text();
        authStatus.textContent = 'Login failed';
        console.error('Login failed:', resp.status, text);
        return;
    }

    const data = await resp.json();
    accessToken = data.access_token || null;

    if (!accessToken) {
        authStatus.textContent = 'Login failed';
        return;
    }

    authStatus.textContent = '✅ Logged in';
}

function loadRooms() {
    fetch('/api/rooms')
        .then(response => response.json())
        .then(data => {
            roomSelect.innerHTML = '';
            (data.rooms || []).forEach(room => {
                const option = document.createElement('option');
                option.value = room;
                option.textContent = room;
                roomSelect.appendChild(option);
            });
        })
        .catch(error => {
            console.error('Error loading rooms:', error);
        });
}

function connectRoom() {
    const room = roomSelect.value || 'general';
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;

    if (!accessToken) {
        alert('Please login first');
        return;
    }

    if (ws) {
        ws.close();
    }

    const wsUrl = `${protocol}//${host}/ws/${room}?token=${encodeURIComponent(accessToken)}`;
    ws = new WebSocket(wsUrl);

    ws.onopen = function() {
        fetch(`/api/messages/${room}`)
            .then(response => response.json())
            .then(data => {
                messagesDiv.innerHTML = '';
                const history = data.messages || [];
                if (history.length > 0) {
                    history.reverse().forEach(msg => {
                        const msgDiv = document.createElement('div');
                        msgDiv.className = 'message';
                        msgDiv.textContent = `${msg.user}: ${msg.text}`;
                        messagesDiv.appendChild(msgDiv);
                    });
                }
                const connectMsg = document.createElement('div');
                connectMsg.className = 'message';
                connectMsg.textContent = `✅ Connected to ${room}`;
                messagesDiv.appendChild(connectMsg);
                messagesDiv.scrollTop = messagesDiv.scrollHeight;
            })
            .catch(error => {
                console.error('History error', error);
            });
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
        messagesDiv.innerHTML += '<div class="message">⛔ Disconnected</div>';
    };
}

function sendMessage() {
    const text = messageInput.value.trim();

    if (!text || !ws || ws.readyState !== WebSocket.OPEN) {
        alert('Not connected or empty message');
        return;
    }

    ws.send(JSON.stringify({ text }));
    messageInput.value = '';
}

loadRooms();
