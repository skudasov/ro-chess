import {joinMsg} from "./msg/msg";
import {blobToText} from "../util/convert";
import {GS} from './game';

const connect = (player) => {
    this.ws = new WebSocket('ws://127.0.0.1:3653');

    this.ws.onopen = () => {
        if (!GS.connected) {
            console.log(`connecting as ${player}`);
            this.ws.send(JSON.stringify(joinMsg(player, player)));
            console.log('connected');
            GS.connected = true;
        }
    };

    this.ws.onmessage = async (event) => {
        let res = await blobToText(event.data);
        console.log(`res: ${res}`);
    };

};
window.connectToServer = connect;

export class GameState {
    constructor() {
        this.connected = false;
        this.figurePool = [];
    }
}