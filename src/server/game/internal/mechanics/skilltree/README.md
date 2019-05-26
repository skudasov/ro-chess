#### Experience

Every unit rewards experience based on how much damage it dealt in combat.
```
Exp = PlayerRewardRatio * dmgToPlayers + UnitRewardRatio * dmgToUnits + ( unitsInDeckMultiplier? )
                                                  eq. 2
```
After XP value for next level is reached unit will reborn in battle and partially healed based on class.

Initially every unit placed on board receives first skill and can level it immidiately.
After that successful combo gives different experience values.
```
Level 1: 100xp 1st profession (2 skills)
Level 2: 200xp 2nd profession (4 skills)
Level 3: 400xp 3rd profession (6 skills)
Level 4: 800xp ?
```
```
ATK 3 Warriors = 80 xp
ATK 2 Warriors = 40 xp
DEF 3 Warriors = 60 xp
DEF 3 Warriors = 30 xp

ATK 3 Mages = 80 xp
ATK 2 Mages = 40 xp
DEF 3 Mages = 60 xp
DEF 2 Mages = 40 xp

ATK 3 Theives = 80 xp
ATK 2 Theives = 40 xp
DEF 3 Theives = 60 xp
DEF 2 Theives = 30 xp

ATK 3 Archers = 80 xp
ATK 2 Archers = 40 xp
DEF 3 Archers = 60 xp
DEF 2 Archers = 30 xp

ATK 3 Acolytes = 80 xp
ATK 2 Acolytes = 40 xp
DEF 3 Acolytes = 60 xp
DEF 2 Acolytes = 30 xp

ATK 3 Merchants = 80 xp
ATK 2 Merchants = 40 xp
DEF 3 Merchants = 60 xp
DEF 2 Merchants = 30 xp
                 tab. 1
```

Combining units in straight combos gives player an opportunity to fast leveling key units for your deck build BUT it's easily countered by opposites.

Experience reward can be balanced by ratios, class combo effectiveness can be balanced by tab.1

Defending all day long will reward minimum experience so you cannot defend forever.

#### Skill trees

Skill tree has 6 tiers, for 3 skill on every row, passive or active.

```
Tier 1:
Tier 2:
Tier 3:
Tier 4:
Tier 5:
Tier 6:
```
Skill leveling workflow:
1. Experience level reached.
2. Command for available skill animation was sent.
3. One skill of current unit level can be learned, client sends his choice in any time before unit have been activated.
4. Skill added to skills rotation.
