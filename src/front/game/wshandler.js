let ws = new WebSocket('ws://127.0.0.1:3653');

ws.onopen = () => {
    console.log('connected');
    let joinMsg = {
        'Join': {
            'Token': 'p1',
            'Name': 'p1',
        }
    };
    ws.send(JSON.stringify(joinMsg));
};

ws.onmessage = (data) => {
    console.log(data)
};
