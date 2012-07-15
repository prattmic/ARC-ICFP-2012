package icfp

import (
        "strconv"
)

/*****************Robot****************/
func (c RobotCell) Parse(coord Coord, mine *Mine) {
    mine.Robot.Coord = coord
}

func (c RobotCell) Update(coord Coord, mine *Mine, updated Map) {
    updated[coord[0]][coord[1]] = mine.Layout[coord[0]][coord[1]]
}

func (c RobotCell) MergeRobot(coord Coord, mine *Mine) bool {
    if mine.Layout[coord[0]][coord[1]] != ROBOT {
        if mine.Layout[coord[0]][coord[1]].MergeRobot(coord, mine) {
            mine.Layout[mine.Robot.Coord[0]][mine.Robot.Coord[1]] = EMPTY
            mine.Layout[coord[0]][coord[1]] = ROBOT
            mine.Robot.Coord = coord
        }
    }

    return true
}

func (c RobotCell) Byte() byte {
    return byte(c)
}

/*****************Rock****************/

func (c RockCell) Parse(coord Coord, mine *Mine) {
    return
}

func (c RockCell) Update(coord Coord, mine *Mine, updated Map) {
    i := coord[0]
    j := coord[1]

    switch {
    case mine.Layout[i+1][j] == EMPTY:
        //Rule 1
        updated[i][j] = EMPTY
        updated[i+1][j] = ROCK

    case (mine.Layout[i+1][j] == ROCK || mine.Layout[i+1][j] == LAMBDA) && mine.Layout[i][j+1] == EMPTY &&  mine.Layout[i+1][j+1] == EMPTY:
        //Rule 2 and 4
        updated[i][j] = EMPTY
        updated[i+1][j+1] = ROCK
    case mine.Layout[i+1][j] == ROCK && mine.Layout[i][j-1] == EMPTY && mine.Layout[i+1][j-1] == EMPTY:
        //Rule 3
        updated[i][j] = EMPTY
        updated[i+1][j-1] = ROCK
    default:
        updated[i][j] = ROCK
    }
}

func (c RockCell) MergeRobot(coord Coord, mine *Mine) bool {
    switch {
    case mine.Robot.Coord[1]<coord[1]:
        mine.Layout[coord[0]][coord[1]+1] = ROCK
        return true
    case mine.Robot.Coord[1]>coord[1]:
        mine.Layout[coord[0]][coord[1]-1] = ROCK
        return true
    }

    return false
}

func (c RockCell) Byte() byte {
    return byte(c)
}

/*****************Wall****************/
func (c WallCell) Parse(coord Coord, mine *Mine) {
    return
}

func (c WallCell) Update(coord Coord, mine *Mine, updated Map) {
    updated[coord[0]][coord[1]] = mine.Layout[coord[0]][coord[1]]
}

func (c WallCell) MergeRobot(coord Coord, mine *Mine) bool {
    return false
}

func (c WallCell) Byte() byte {
    return byte(c)
}

/*****************Lambda****************/
func (c LambdaCell) Parse(coord Coord, mine *Mine) {
    mine.Lambda = append(mine.Lambda, coord)
}

func (c LambdaCell) Update(coord Coord, mine *Mine, updated Map) {
    updated[coord[0]][coord[1]] = mine.Layout[coord[0]][coord[1]]
}

func (c LambdaCell) MergeRobot(coord Coord, mine *Mine) bool {
    err := mine.eatLambda(coord)
    if err != nil {
        return false
    }
    return true
}

func (c LambdaCell) Byte() byte {
    return byte(c)
}

/*****************Earth****************/
func (c EarthCell) Parse(coord Coord, mine *Mine) {
    return
}

func (c EarthCell) Update(coord Coord, mine *Mine, updated Map) {
    updated[coord[0]][coord[1]] = mine.Layout[coord[0]][coord[1]]
}

func (c EarthCell) MergeRobot(coord Coord, mine *Mine) bool {
    return true
}

func (c EarthCell) Byte() byte {
    return byte(c)
}

/*****************Empty****************/
func (c EmptyCell) Parse(coord Coord, mine *Mine) {
    return
}

func (c EmptyCell) Update(coord Coord, mine *Mine, updated Map) {
    if updated[coord[0]][coord[1]] != ROCK && updated[coord[0]][coord[1]] != BEARD {
        updated[coord[0]][coord[1]] = EMPTY
    }
}

func (c EmptyCell) MergeRobot(coord Coord, mine *Mine) bool {
    return true
}

func (c EmptyCell) Byte() byte {
    return byte(c)
}

/*****************Lift****************/
func (c LiftCell) Parse(coord Coord, mine *Mine) {
    mine.Lift.Coord = coord
    mine.Lift.Open = false
}

func (c LiftCell) Update(coord Coord, mine *Mine, updated Map) {
    updated[coord[0]][coord[1]] = mine.Layout[coord[0]][coord[1]]
}

func (c LiftCell) MergeRobot(coord Coord, mine *Mine) bool {
    if c == OLIFT {
        mine.Complete = true
        return true
    }

    return false
}

func (c LiftCell) Byte() byte {
    return byte(c)
}

/*****************Tramp****************/
func (c TrampCell) Parse(coord Coord, mine *Mine) {
    targ := mine.Trampolines[string(mine.Layout[coord[0]][coord[1]].Byte())]
    targ.TrampCoord = coord
    mine.Trampolines[string(mine.Layout[coord[0]][coord[1]].Byte())] = targ
}

func (c TrampCell) Update(coord Coord, mine *Mine, updated Map) {
    updated[coord[0]][coord[1]] = mine.Layout[coord[0]][coord[1]]
}

func (c TrampCell) MergeRobot(coord Coord, mine *Mine) bool {
    mine.takejump(coord)
    // Note in this special case we always return false because takejump handles the robot's movement
    return false    
}

func (c TrampCell) Byte() byte {
    return byte(c)
}

/*****************Targ****************/
func (c TargCell) Parse(coord Coord, mine *Mine) {
    num, _ := strconv.Atoi(string(mine.Layout[coord[0]][coord[1]].Byte()))
    mine.TargetCoord(num, coord)
}

func (c TargCell) Update(coord Coord, mine *Mine, updated Map) {
    updated[coord[0]][coord[1]] = mine.Layout[coord[0]][coord[1]]
}

func (c TargCell) MergeRobot(coord Coord, mine *Mine) bool {
    return false
}

func (c TargCell) Byte() byte {
    return byte(c)
}

/*****************Beard****************/
func (c BeardCell) Parse(coord Coord, mine *Mine) {
    return
}

func (c BeardCell) Update(coord Coord, mine *Mine, updated Map) {
    i := coord[0]
    j := coord[1]

    updated[i][j] = mine.Layout[i][j]
    if mine.Gcount == 0 {
        for k := i-1; k <= i+1; k++ {
            for l := j-1; l <= j+1; l++ {
                if mine.Layout[k][l] == EMPTY {
                    updated[k][l] = BEARD
                }
            }
        }
    }
}

func (c BeardCell) MergeRobot(coord Coord, mine *Mine) bool {
    return false
}

func (c BeardCell) Byte() byte {
    return byte(c)
}

/*****************Razor****************/
func (c RazorCell) Parse(coord Coord, mine *Mine) {
    return
}

func (c RazorCell) Update(coord Coord, mine *Mine, updated Map) {
    updated[coord[0]][coord[1]] = mine.Layout[coord[0]][coord[1]]
}

func (c RazorCell) MergeRobot(coord Coord, mine *Mine) bool {
    mine.Robot.Razors++
    return true
}

func (c RazorCell) Byte() byte {
    return byte(c)
}
