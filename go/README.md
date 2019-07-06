# 6

### Strategy
- Using __Alpha-Beta pruning__ and expect next situations with __depth 6__
(it should be changed dinamically with the elapsed time, but now it is fixed).
- Scores of each situation are get from the score table.
I got it from here:
https://nlab.itmedia.co.jp/nl/articles/1710/06/news013_2.html

![othello_score](https://user-images.githubusercontent.com/34668695/60395894-165a3400-9b75-11e9-9a4c-29fb915e58d4.jpg)

- → Changed the score table to the one from here
(https://uguisu.skr.jp/othello/5-1.html)
and adjusted it by myself.
<img width="227" alt="Screen Shot 2019-07-06 at 13 30 07" src="https://user-images.githubusercontent.com/34668695/60751572-5a748b00-9ff2-11e9-83b6-eb8e5b4593d4.png">


- The __score table is different__ slightly __between first and middle phase__.

- In __last phase__, the strategy changes 
from getting most scores to __getting most pieces__ 
because the program can estimate all of the possible situation.

- A parameter from __the number of valid moves for the next player__ is added to the score
before last phase.  