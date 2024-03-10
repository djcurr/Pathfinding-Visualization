// src/app/grid/grid.component.ts
import {
  Component,
  OnChanges,
  Input,
  OnInit,
  OnDestroy,
  SimpleChanges,
} from '@angular/core';
import { Cell } from '../cell/cell.model';
import {WasmService, Point} from "../wasm.service";
import {DefaultAnimationSpeed, DefaultGridSize} from "../app.component";

@Component({
  selector: 'app-grid',
  templateUrl: './grid.component.html',
  styleUrls: ['./grid.component.css']
})
export class GridComponent implements OnInit, OnDestroy, OnChanges {
  grid: Cell[][] = [];
  isMouseDown: boolean = false;
  startSelected: boolean = false;
  endSelected: boolean = false;
  prevCell: Cell = new Cell(0, 0);
  solved: boolean = false;
  @Input() isEraserActive: boolean = false;
  @Input() animationSpeed: number = DefaultAnimationSpeed;
  @Input() gridSize: number = DefaultGridSize;
  @Input() clearGridEvent!: any;
  @Input() startPathfinding!: boolean;
  @Input() drawGridEvent!: boolean;


  constructor(private wasmService: WasmService) { }

  ngOnInit(): void {
    this.drawGrid()
    document.addEventListener('mouseup', this.onGlobalMouseUp.bind(this));
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes["gridSize"] && changes["gridSize"].previousValue !== changes["gridSize"].currentValue) {
      this.changeGridSize(changes["gridSize"].currentValue);
      requestAnimationFrame(this.drawGrid.bind(this));
    } else if (changes["clearGridEvent"]) {
      this.clearGrid();
    } else if (changes["startPathfinding"]) {
      this.solveGrid();
    } else if (changes["drawGridEvent"]) {
      if (!this.solved) {
        requestAnimationFrame(this.drawGrid.bind(this));
      }
    }
  }

  ngOnDestroy(): void {
    document.removeEventListener('mouseup', this.onGlobalMouseUp.bind(this));
  }

  onGlobalMouseUp(): void {
    this.isMouseDown = false;
    this.startSelected = false;
    this.endSelected = false;
  }

  onMouseDown(cell: Cell): void {
    this.isMouseDown = true;
    if (this.prevCell.x === cell.x && this.prevCell.y === cell.y || this.solved) { return; }
    if (cell.isStart) {
      this.startSelected = true;
    } else if (cell.isEnd) {
      this.endSelected = true;
    } else {
      this.wasmService.setWall(cell.x, cell.y, !this.isEraserActive).catch((error) => {
        console.error("Error setting wall:", error);
      });
      this.grid[cell.y][cell.x] = {...this.grid[cell.y][cell.x], isWall: !this.isEraserActive};
    }
    this.prevCell = cell;
  }

  onMouseEnter(cell: Cell): void {
    if (this.prevCell.x === cell.x && this.prevCell.y === cell.y || this.solved) { return; }
    if (this.isMouseDown) {
      if (this.startSelected && !cell.isWall && !cell.isEnd) {
        this.wasmService.setStart(cell.x, cell.y).catch((error) => {
          console.error("Error setting start:", error);
        }).then(() => {
          requestAnimationFrame(this.drawGrid.bind(this));
        })
      } else if (this.endSelected && !cell.isWall && !cell.isStart) {
        this.wasmService.setEnd(cell.x, cell.y).catch((error) => {
          console.error("Error setting end:", error);
        }).then(() => {
          requestAnimationFrame(this.drawGrid.bind(this));
        })
      } else {
        this.wasmService.setWall(cell.x, cell.y, !this.isEraserActive).catch((error) => {
          console.error("Error setting wall:", error);
        })
        this.grid[cell.y][cell.x] = {...this.grid[cell.y][cell.x], isWall: !this.isEraserActive};
      }
      this.prevCell = cell;
    }
  }

  onMouseUp(): void {
    this.isMouseDown = false;
    this.endSelected = false;
    this.startSelected = false;
  }

  solveGrid(): void {
    this.wasmService.findPath().then(() => {
      this.solved = true;
      requestAnimationFrame(this.showSnapshots.bind(this));
    }).catch((error) => {
      console.error("Error finding path:", error);
    });
  }

  drawGrid(): void {
    this.wasmService.getGrid().then((result) => {
      this.grid = result
    }, (error) => {
      console.error("Error finding grid:", error);
    })
  }

  changeGridSize(size: number): void {
    this.wasmService.changeGridSize(size, size).catch((error) => {
      console.error("Error setting size:", error);
    });
    this.solved = false;
    requestAnimationFrame(this.drawGrid.bind(this));
  }

  clearGrid(): void {
    this.wasmService.clearGrid().catch((error) => {
      console.error("Error clearing grid:", error);
    })
    this.solved = false;
    requestAnimationFrame(this.drawGrid.bind(this));
  }

  showSnapshots(): void {
    this.wasmService.getSnapshot().then((result) => {
      if (result.length === 0) {
        requestAnimationFrame(this.showPath.bind(this))
        return
      }
      this.grid = result;
      const delay = 300 / this.animationSpeed
      setTimeout(() => requestAnimationFrame(this.showSnapshots.bind(this)), delay)
    })
  }

  async showPath(): Promise<void> {
    try {
      const result = await this.wasmService.getPath();
      const start = await this.wasmService.getStart();
      if (!start) {
        console.error("Error getting start");
        return;
      }

      const end = await this.wasmService.getEnd();
      if (!end) {
        console.error("Error getting end");
        return;
      }
      requestAnimationFrame(this.reconstructPath.bind(this, result, start, result.get(end.toString())));
    } catch (error) {
      console.error("An error occurred:", error);
    }
  }

  reconstructPath(path: Map<String, Point>, goal: Point, current: Point | undefined) {
    if (!current) {
      return
    } else {
      this.grid[current.getY()][current.getX()] = {...this.grid[current.getY()][current.getX()], isPath: true, visited: false};
      const next = path.get(current.toString())
      if (!next || (next.getX() === goal.getX() && next.getY() === goal.getY())) {
        return
      }
      this.reconstructPath(path, goal, next);
    }
  }
}
