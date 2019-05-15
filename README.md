##  RO chess

A card/board game, that mixes strategy, economy and tactics, partially based on Ragnarok Online universe.

### Prototype concepts

Battlefield consists of three zones (w x h):
- P1 start zone (9x4) 
- P2 start zone (9x4) 
- Neutral zone (9x5)

The key moment here is that player sees only his start zone in the beginning, reconnaissance is the key to victory or solid defence.

```
[---------]
[---------]
[---------]
[---------] <- P1 start zone
<--------->
<--------->
<---------> <- neutral creeps line
<--------->
<--------->
[---------] <- P2 start zone
[---------]
[---------]
[---------]
```

Before a match players can pick units they like, of one side (orcs/humans) and add it to deck.
When the match started, players can ban each other deck units (up to 2 for now).

After that the game begins.

Main resource is gold, it' generated through turns slowly, can be mined with peon/peasant combos.

Every player has 100 HP, 0 MP, 0 gold and no buffs at the beginning.
Each player receives a pool of units in his turn (5 units), from which he can choose only 1 unit to move on the board.
Player can lock the pool rotation, but this lock is limited by 5 turns (you should build deck carefully if you need predictability, it's part of a strategy).

After player puts the figure on the board, he can activate any amount of figures in one turn.
When activated figures will march front until enemy scoring zone or enemy unit/trap/spell/defence reached.
The goal of a game is not to spam units without strategy, singular unit activation is limited to 2 units simultaneously scouting.
Every unit that passes enemy scoring zone will attack enemy player using simple formula:
```
dmgToPlayer = unitHP * unitAvgDps
                            eq. 1
```

Main objective of the game is to kill opponent by applying different unit combos, with corresponding knowledge of opponents moves/units, e.g.:
```
A
[---------]
[-G-------]
[-G-------]
[-G-------] <- gives grunts attack buff and makes a squad

B
[---------]
[---------]
[---------]
[-GGG-----] <- gives armor buff and makes a squad, but making them unable to move
```
You can activate such a squad (A), and it will move as singular unit towards opponent, some squad bonuses can be lost if one squad member dies.
Defensive squads (B) cannot move, but can improve their bonuses over time slightly.

All units, even in squads fight one by one in turn with buffs applied.
When squad reaches opponent score zone their do dmg according to eq.1, overall dps can be changed with squad buffs.

### Reconnaissance mechanics

*Under construction*

### Neutral creeps mechanics

*Under construction*

## Development

### OS X

To run lint/test/server + front dev server:

```
make all
```
For other commands see Makefile

### Global TODO's
- [x] Refactor FixedTurns crutch and make multiboard connections
- [x] Make new type inside UpdateBatch msg, combat log for front rendering
- [ ] Test combo stacks, test them simultaneously
- [ ] Make activation command to be applied for all units in column

