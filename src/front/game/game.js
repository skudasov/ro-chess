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


// var game = new Phaser.Game(800, 600, Phaser.CANVAS, 'phaser-example', { preload: preload, create: create, update: update, render: render });
//
// // window.game = game;
//
// function preload() {
//
//     game.load.tilemap('matching', 'assets/tilemaps/maps/phaser_tiles.json', null, Phaser.Tilemap.TILED_JSON);
//     game.load.image('tiles', 'assets/tilemaps/tiles/phaser_tiles.png');//, 100, 100, -1, 1, 1);
// }
//
// var timeCheck = 0;
// var flipFlag = false;
//
// var startList = new Array();
// var squareList = new Array();
//
// var masterCounter = 0;
// var squareCounter = 0;
// var clickCount = 0;
// var square1Num;
// var square2Num;
// var savedSquareX1;
// var savedSquareY1;
// var savedSquareX2;
// var savedSquareY2;
//
// var map;
// var tileset;
// var layer;
//
// var marker;
// var currentTile;
// var currentTilePosition;
//
// var tileBack = 25;
// var timesUp = '+';
// var youWin = '+';
//
// var myCountdownSeconds;
//
//
// function create() {
//
//     map = game.add.tilemap('matching');
//
//     map.addTilesetImage('Desert', 'tiles');
//
//     //tileset = game.add.tileset('tiles');
//
//     layer = map.createLayer('Ground');//.tilemapLayer(0, 0, 600, 600, tileset, map, 0);
//
//     //layer.resizeWorld();
//
//     marker = game.add.graphics();
//     marker.lineStyle(2, 0x00FF00, 1);
//     marker.drawRect(0, 0, 100, 100);
//
//     randomizeTiles();
//
// }
//
// function update() {
//
//     countDownTimer();
//
//     if (layer.getTileX(game.input.activePointer.worldX) <= 5) // to prevent the marker from going out of bounds
//     {
//         marker.x = layer.getTileX(game.input.activePointer.worldX) * 100;
//         marker.y = layer.getTileY(game.input.activePointer.worldY) * 100;
//     }
//
//     if (flipFlag == true)
//     {
//         if (game.time.totalElapsedSeconds() - timeCheck > 0.5)
//         {
//             flipBack();
//         }
//     }
//     else
//     {
//         processClick();
//     }
// }
//
//
// function countDownTimer() {
//
//     var timeLimit = 120;
//
//     mySeconds = game.time.totalElapsedSeconds();
//     myCountdownSeconds = timeLimit - mySeconds;
//
//     if (myCountdownSeconds <= 0)
//     {
//         // time is up
//         timesUp = 'Time is up!';
//         myCountdownSeconds = 0;
//
//     }
// }
//
// function processClick() {
//
//     currentTile = map.getTile(layer.getTileX(marker.x), layer.getTileY(marker.y));
//     currentTilePosition = ((layer.getTileY(game.input.activePointer.worldY)+1)*6)-(6-(layer.getTileX(game.input.activePointer.worldX)+1));
//
//     if (game.input.mousePointer.isDown)
//     {
//         // check to make sure the tile is not already flipped
//         if (currentTile.index == tileBack)
//         {
//             // get the corresponding item out of squareList
//             currentNum = squareList[currentTilePosition-1];
//             flipOver();
//             squareCounter++;
//             clickCount++;
//
//             // is the second tile of pair flipped?
//             if  (squareCounter == 2)
//             {
//                 // reset squareCounter
//                 squareCounter = 0;
//                 square2Num = currentNum;
//                 // check for match
//                 if (square1Num == square2Num)
//                 {
//                     masterCounter++;
//
//                     if (masterCounter == 18)
//                     {
//                         // go "win"
//                         youWin = 'Got them all!';
//                         if (clickCount == 18)
//                         {
//                             youWin = 'Hard-mode achieved';
//                         }
//                     }
//                     else
//                     {
//                         savedSquareX2 = layer.getTileX(marker.x);
//                         savedSquareY2 = layer.getTileY(marker.y);
//                         flipFlag = true;
//                         timeCheck = game.time.totalElapsedSeconds();
//                     }
//                 }
//                 else
//                 {
//                     savedSquareX1 = layer.getTileX(marker.x);
//                     savedSquareY1 = layer.getTileY(marker.y);
//                     square1Num = currentNum;
//                 }
//             }
//         }
//     }
// }
//
// function flipOver() {
//
//     map.putTile(currentNum, layer.getTileX(marker.x), layer.getTileY(marker.y));
// }
//
// function flipBack() {
//
//     flipFlag = false;
//
//     map.putTile(tileBack, savedSquareX1, savedSquareY1);
//     map.putTile(tileBack, savedSquareX2, savedSquareY2);
//
// }
//
// function randomizeTiles() {
//
//     for (num = 1; num <= 18; num++)
//     {
//         startList.push(num);
//     }
//     for (num = 1; num <= 18; num++)
//     {
//         startList.push(num);
//     }
//
//     // for debugging
//     myString1 = startList.toString();
//
//     // randomize squareList
//     for (i = 1; i <=36; i++)
//     {
//         var randomPosition = game.rnd.integerInRange(0,startList.length - 1);
//
//         var thisNumber = startList[ randomPosition ];
//
//         squareList.push(thisNumber);
//         var a = startList.indexOf(thisNumber);
//
//         startList.splice( a, 1);
//     }
//
//     // for debugging
//     myString2 = squareList.toString();
//
//     for (col = 0; col < 6; col++)
//     {
//         for (row = 0; row < 6; row++)
//         {
//             map.putTile(tileBack, col, row);
//         }
//     }
// }
//
// function getHiddenTile() {
//
//     thisTile = squareList[currentTilePosition-1];
//     return thisTile;
// }
//
// function render() {
//
//     game.debug.text(timesUp, 620, 208, 'rgb(0,255,0)');
//     game.debug.text(youWin, 620, 240, 'rgb(0,255,0)');
//
//     game.debug.text('Time: ' + myCountdownSeconds, 620, 15, 'rgb(0,255,0)');
//
//     //game.debug.text('squareCounter: ' + squareCounter, 620, 272, 'rgb(0,0,255)');
//     game.debug.text('Matched Pairs: ' + masterCounter, 620, 304, 'rgb(0,0,255)');
//     game.debug.text('Matched Pairs: ' + clickCount, 620, 320, 'rgb(0,0,255)');
//
//
//     //game.debug.text('startList: ' + myString1, 620, 208, 'rgb(255,0,0)');
//     //game.debug.text('squareList: ' + myString2, 620, 240, 'rgb(255,0,0)');
//
//
//     game.debug.text('Tile: ' + map.getTile(layer.getTileX(marker.x), layer.getTileY(marker.y)).index, 620, 48, 'rgb(255,0,0)');
//
//     game.debug.text('LayerX: ' + layer.getTileX(marker.x), 620, 80, 'rgb(255,0,0)');
//     game.debug.text('LayerY: ' + layer.getTileY(marker.y), 620, 112, 'rgb(255,0,0)');
//
//     game.debug.text('Tile Position: ' + currentTilePosition, 620, 144, 'rgb(255,0,0)');
//     game.debug.text('Hidden Tile: ' + getHiddenTile(), 620, 176, 'rgb(255,0,0)');
// }