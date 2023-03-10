import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-modal-input',
  templateUrl: './modal-input.component.html',
  styleUrls: ['./modal-input.component.scss']
})
export class ModalInputComponent {
    inputTitle: string | undefined = "Nickname" 

    @Output() inputText = new EventEmitter<string>()
    @Input() active: boolean = false 
    @Output() onConfirm: EventEmitter<any> = new EventEmitter() 

    confirmClicked(): void {
        this.onConfirm.emit()
    }

    updateInputText (event: Event): void {
        this.inputText.emit((event.target as HTMLInputElement).value)
    }

    closeModal(): void {
        this.active = false
    }
}
