import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { ProductService } from '../../../services/product.service';
import { CreateProductDTO } from '../../../models/product.model';

@Component({
  selector: 'app-product-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatCardModule,
    MatIconModule,
    MatSnackBarModule
  ],
  templateUrl: './product-form.component.html',
  styleUrls: ['./product-form.component.scss']
})
export class ProductFormComponent implements OnInit {
  productForm: FormGroup;
  loading = false;

  constructor(
    private fb: FormBuilder,
    private productService: ProductService,
    private router: Router,
    private snackBar: MatSnackBar
  ) {
    // Inicializa o formulário com validações
    this.productForm = this.fb.group({
      code: ['', [Validators.required, Validators.minLength(3)]],
      description: ['', [Validators.required, Validators.minLength(5)]],
      balance: [0, [Validators.required, Validators.min(0)]]
    });
  }

  ngOnInit(): void {}

  /**
   * Submete o formulário
   * Usa switchMap do RxJS para encadear observables
   */
  onSubmit(): void {
    if (this.productForm.invalid) {
      this.showError('Por favor, preencha todos os campos corretamente');
      return;
    }

    this.loading = true;
    const product: CreateProductDTO = this.productForm.value;

    this.productService.createProduct(product).subscribe({
      next: () => {
        this.showSuccess('Produto cadastrado com sucesso!');
        this.router.navigate(['/products']);
      },
      error: (error) => {
        this.showError(error.message);
        this.loading = false;
      }
    });
  }

  /**
   * Cancela e volta para a lista
   */
  onCancel(): void {
    this.router.navigate(['/products']);
  }

  /**
   * Mostra mensagem de sucesso
   */
  private showSuccess(message: string): void {
    this.snackBar.open(message, 'Fechar', {
      duration: 3000,
      horizontalPosition: 'end',
      verticalPosition: 'top',
      panelClass: ['success-snackbar']
    });
  }

  /**
   * Mostra mensagem de erro
   */
  private showError(message: string): void {
    this.snackBar.open(message, 'Fechar', {
      duration: 5000,
      horizontalPosition: 'end',
      verticalPosition: 'top',
      panelClass: ['error-snackbar']
    });
  }

  /**
   * Getters para facilitar acesso aos controles do formulário
   */
  get code() {
    return this.productForm.get('code');
  }

  get description() {
    return this.productForm.get('description');
  }

  get balance() {
    return this.productForm.get('balance');
  }
}