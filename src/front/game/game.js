import config from '../config'
import {GameState} from "./state";

let map;
let boardLayer;
let boardLayerName = 'board';
let bgLayer;
let bgLayerName = 'background';
let poolLayer;
let poolLayerName = 'pool';
let figuresLayer;
let figuresLayerName = 'figures';

let tileNumFigure = 196;

export let GS = new GameState();

export default class w2chGame extends Phaser.State {
    init(args) {
    }

    preload(game) {
        game.load.tilemap(config.tileMapName, config.tileMapPath, null, Phaser.Tilemap.TILED_JSON);
        game.load.image('tiles-key', '../assets/maps/terrain.png');

        game.load.image('tower', '../assets/sprites/tower-32.png');
        game.load.image('bullet', '../assets/sprites/bullet.png');
        game.load.spritesheet('duck', '../assets/sprites/duck.png', 32, 32, 8);
        game.load.spritesheet('panda', '../assets/sprites/panda.png', 32, 32, 3);
        game.load.spritesheet('dog', '../assets/sprites/dog.png', 32, 32, 6);
        game.load.spritesheet('penguin', '../assets/sprites/penguin.png', 32, 32, 4);
    }

    create(game) {
        game.physics.startSystem(Phaser.Physics.ARCADE);
        map = game.add.tilemap(config.tileMapName);

        map.addTilesetImage('terrain', 'tiles-key');
        bgLayer = map.createLayer(bgLayerName);
        boardLayer = map.createLayer(boardLayerName);
        poolLayer = map.createLayer(poolLayerName);
        figuresLayer = map.createLayer(figuresLayerName);
        fillPoolWithFigures(game)
    }

    update(game) {
        processClick(game)
    }
}

const fillPoolWithFigures = (game) => {
    let startTileNum = tileNumFigure;
    let endTileNum = startTileNum + 6;
    let tileX = 2;
    let tileY = 12;
  for (let i = tileNumFigure; i < endTileNum; i++) {
      console.log(`putting tile on ${tileX}, ${tileY}`);
      map.putTile(tileNumFigure, tileX, tileY, poolLayerName);
      tileY++;
  }
};

const processClick = (game) => {
    let x = game.input.activePointer.worldX;
    let y = game.input.activePointer.worldY;
    if (game.input.mousePointer.isDown) {
        console.log(`clicking on ${x}, ${y}`);
        map.putTile(tileNumFigure, boardLayer.getTileX(x), boardLayer.getTileY(y), figuresLayerName);
    }
    let currentTile = map.getTile(boardLayer.getTileX(x), boardLayer.getTileY(y));
    console.log(`current tile: ${currentTile}`);
};