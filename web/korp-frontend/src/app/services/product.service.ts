import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError, BehaviorSubject } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { Product, CreateProductDTO, UpdateProductDTO } from '../models/product.model';

@Injectable({
  providedIn: 'root'
})
export class ProductService {
  private apiUrl = 'http://localhost:8081/api/products';
  
  // BehaviorSubject para manter lista de produtos em memória (cache)
  private productsSubject = new BehaviorSubject<Product[]>([]);
  public products$ = this.productsSubject.asObservable();

  constructor(private http: HttpClient) {}

  /**
   * Lista todos os produtos
   * Usa RxJS operators: tap para atualizar cache, catchError para tratamento
   */
  getProducts(): Observable<Product[]> {
    return this.http.get<Product[]>(this.apiUrl).pipe(
      tap(products => this.productsSubject.next(products)),
      catchError(this.handleError)
    );
  }

  /**
   * Busca um produto por ID
   */
  getProduct(id: string): Observable<Product> {
    return this.http.get<Product>(`${this.apiUrl}/${id}`).pipe(
      catchError(this.handleError)
    );
  }

  /**
   * Cria um novo produto
   */
  createProduct(product: CreateProductDTO): Observable<Product> {
    return this.http.post<Product>(this.apiUrl, product).pipe(
      tap(() => this.getProducts().subscribe()), // Recarrega lista
      catchError(this.handleError)
    );
  }

  /**
   * Atualiza um produto existente
   */
  updateProduct(id: string, product: UpdateProductDTO): Observable<Product> {
    return this.http.put<Product>(`${this.apiUrl}/${id}`, product).pipe(
      tap(() => this.getProducts().subscribe()), // Recarrega lista
      catchError(this.handleError)
    );
  }

  /**
   * Tratamento centralizado de erros
   * Retorna mensagens amigáveis para o usuário
   */
  private handleError(error: HttpErrorResponse) {
    let errorMessage = 'Ocorreu um erro desconhecido';
    
    if (error.error instanceof ErrorEvent) {
      // Erro do lado do cliente
      errorMessage = `Erro: ${error.error.message}`;
    } else {
      // Erro do lado do servidor
      if (error.status === 0) {
        errorMessage = 'Não foi possível conectar ao servidor. Verifique se o Stock Service está rodando.';
      } else if (error.status === 404) {
        errorMessage = 'Produto não encontrado';
      } else if (error.status === 409) {
        errorMessage = 'Código de produto já existe';
      } else if (error.error?.message) {
        errorMessage = error.error.message;
      } else {
        errorMessage = `Erro do servidor: ${error.status}`;
      }
    }
    
    console.error('Erro no ProductService:', error);
    return throwError(() => new Error(errorMessage));
  }
}