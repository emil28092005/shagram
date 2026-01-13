let ws = null;
const messagesDiv = document.getElementById('messages');
const messageInput = document.getElementById('messageInput');
const roomSelect = document.getElementById('roomSelect');

function loadRooms() {
    fetch('/api/rooms')
    .then(response => response.json())
    .then(data => {
        roomSelect.innerHTML = '';
        data.rooms.forEach(room => {
            const option = document.createElement('option');
            option.value = room;
            option.textContent = room;
            roomSelect.appendChild(option);
        });
    })
    .catch(error => {
        console.error('Error loading rooms:', error);
    })
}

function connectRoom() {
    const room = roomSelect.value || 'general';
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    
    if (ws) {
        ws.close();
    }

    ws = new WebSocket(`${protocol}//localhost:8080/ws/${room}`);

    

    ws.onopen = function() {
        console.log('WebSocket connected, loading history for', room);

        fetch(`/api/messages/${room}`)
            .then(response => {
                console.log('API response:', response.status);
                return response.json()
            })
            .then(data => {
                console.log('History data:', data);
                messagesDiv.innerHTML = '';
                const history = data.messages || [];
                console.log('history length:', history.length);
                if (history.length > 0) {
                    history.reverse().forEach((msg, index) => {
                        console.log(`Message ${index}:`, msg)
                        const msgDiv = document.createElement('div');
                        msgDiv.className = 'message';
                        msgDiv.textContent = `${msg.user}: ${msg.text}`;
                        messagesDiv.appendChild(msgDiv);
                    })
                }
                const connectMsg = document.createElement('div');
                connectMsg.className = 'message';
                connectMsg.textContent = `✅ Connected to ${room}`;
                messagesDiv.appendChild(connectMsg);
                messagesDiv.scrollTop = messagesDiv.scrollHeight;
                console.log('History loaded!');
            })
            .catch(error => {
                console.error('❌ History error', error);
            })

        
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
    const usernameInput = document.getElementById('userInput');
    const username = usernameInput ? usernameInput.value.trim() || 'Anonymous' : 'Anonymous';

    if (!text || !ws || ws.readyState !== WebSocket.OPEN) {
        alert('Not connected or empty message');
        return;
    }
    
    ws.send(JSON.stringify({
        text: text,
        user: username,
        }));
    messageInput.value = '';
}

loadRooms();