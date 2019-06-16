import {joinMsg} from "./msg/msg";

let ws = new WebSocket('ws://127.0.0.1:3653');

ws.onopen = () => {
    console.log('connecting');
    ws.send(JSON.stringify(joinMsg('p1', 'p1')));
    console.log('connected');
};

ws.onmessage = async (event) => {
    let res = await blobToText(event.data);
    console.log(`res: ${res}`);
};

function blobToText(file){
    return new Promise((resolve, reject) => {
        let fr = new FileReader();
        fr.onload = () => {
            resolve(fr.result)
        };
        fr.readAsText(file);
    });
}
