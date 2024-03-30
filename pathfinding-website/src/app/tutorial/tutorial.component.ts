import { Component } from '@angular/core';
import {NgIf} from "@angular/common";
import {NgbModal} from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: 'app-tutorial',
  standalone: true,
  imports: [
    NgIf
  ],
  templateUrl: './tutorial.component.html',
  styleUrl: './tutorial.component.scss'
})

export class TutorialComponent {
  steps = [
    { title: 'Walls', gifUrl: 'assets/walls.gif', description: 'Click and drag on the grid to draw walls.' },
    { title: 'Eraser', gifUrl: 'assets/eraser.gif', description: 'Click the Eraser button to toggle the eraser. Then, click and drag on the grid to erase walls. Clear all walls by clicking clear.' },
    { title: 'Algorithm', gifUrl: 'assets/algorithm.gif', description: 'Click the currently selected pathfinding algorithm to modify it.' },
    { title: 'Size & Speed', gifUrl: 'assets/size.gif', description: 'Change the size of the grid and the speed of the algorithm visualization by changing the respective sliders. Click Reset to reset all values to their defaults.' },
    { title: 'Start & End Nodes', gifUrl: 'assets/move.gif', description: 'Move the start and end nodes by clicking and dragging.' },
    { title: 'Maze', gifUrl: 'assets/maze.gif', description: 'Click Create Maze to generate a random maze. Clear the grid before creating a new maze.' },
    { title: 'Start Pathfinding', gifUrl: 'assets/start.gif', description: 'Click Start Pathfinding to start the algorithm. Clear the grid to try a different configuration.' },
  ];

  constructor(private modalService: NgbModal) {}

  currentStepIndex = 0;

  get currentStep() {
    return this.steps[this.currentStepIndex];
  }

  nextStep() {
    if (this.currentStepIndex < this.steps.length - 1) {
      this.currentStepIndex++;
    } else {
      this.closeModal();
      this.currentStepIndex = 0;
    }
  }

  previousStep() {
    if (this.currentStepIndex > 0) {
      this.currentStepIndex--;
    }
  }

  closeModal() {
    this.modalService.dismissAll();
  }
}
