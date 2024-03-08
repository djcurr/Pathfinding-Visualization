// src/app/models/cell.model.ts
export class Cell {
  constructor(
    public x: number,
    public y: number,
    public isWall: boolean = false,
    public isStart: boolean = false,
    public isEnd: boolean = false,
    public isPath: boolean = false,
    public visited: boolean = false
  ) {}
}
