import {joinMsg, turnEndMsg} from "./msg/msg";
import {blobToText} from "../util/convert";
import {GS} from './game';

export class GameState {
    constructor() {
        this.ws = new WebSocket('ws://127.0.0.1:3653');

        this.ws.onmessage = async (event) => {
            let res = await blobToText(event.data);
            let msg = JSON.parse(res);
            let msgType = Object.keys(msg);
            switch (true) {
                case (msgType == 'Joined'):
                    console.log(`joined msg received`);
                    this.boardName = msg.Joined.Board;
                    break;
                default:
                    console.log(`unknown msg type received: ${msgType}`)
            }
            console.log(`msg: ${JSON.stringify(msg)}`);
        };

        this.connected = false;
        this.player = null;
        this.boardName = null;
        this.figurePool = [];
    }

    connect(player) {
        if (!this.connected) {
            console.log(`connecting as ${player}`);
            this.ws.send(JSON.stringify(joinMsg(player, player)));
            console.log('connected');
            this.connected = true;
            this.player = player;
        }
    }
    // debug
    rewind(turns){
        this.ws.send(JSON.stringify(turnEndMsg(this.player, this.boardName)))
    };
}

export const connect = (player) => {
    GS.connect(player);
};

export const rewindTurns = (turns) => {
    GS.rewind(turns);
};
window.connectToServer = connect;
window.rewindTurns = rewindTurns;