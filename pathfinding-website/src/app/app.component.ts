import { Component } from '@angular/core';
import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'] // Corrected property name and usage
})
export class AppComponent {
  title = 'Pathfinding Visualization';
  isEraserActive: boolean = false;
  gridSize: number = DefaultGridSize;
  animationSpeed: number = DefaultAnimationSpeed;
  clearGrid: boolean = true;
  startPathfinding: boolean = false;
  drawGrid: boolean = false;

  constructor(private breakpointObserver: BreakpointObserver) {
    this.breakpointObserver.observe([
      Breakpoints.XSmall,
      Breakpoints.Small,
      Breakpoints.Medium,
      Breakpoints.Large,
      Breakpoints.XLarge,
    ]).subscribe(result => {
      if (result.matches) {
        if (result.breakpoints[Breakpoints.XSmall] || result.breakpoints[Breakpoints.Small]) {
          this.gridSize = 13;
          DefaultGridSize = 13;
        }
      }
    });
  }

  handleEraserToggled(isEraserActive: boolean): void {
    this.isEraserActive = isEraserActive;
  }
  handleAnimationSpeedChanged(newSpeed: number): void {
    this.animationSpeed = newSpeed;
  }

  handleStartPathfindingEvent() {
    this.startPathfinding = !this.startPathfinding;
  }

  handleChangeGridSize(size: number) {
    this.gridSize = size;
  }

  handleClearGrid() {
    this.clearGrid = !this.clearGrid;
  }

  handleDrawGrid() {
    this.drawGrid = !this.drawGrid;
  }
}

export var DefaultGridSize: number = 31;
export const DefaultAnimationSpeed: number = 20;
export const AStar: number = 0;
export const Dijkstra: number = 1;
export const DefaultAlgorithm: number = AStar;
export const Algorithms: Map<number, string> = new Map([
  [AStar, "A*"],
  [Dijkstra, "Dijkstra"]
])
