import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs'; // Uso de RxJS conforme requisito

@Injectable({
  providedIn: 'root',
})
export class ProdutoService {
  // Porta 8081 dedicada ao Serviço de Estoque
  private apiUrl = 'http://localhost:8081/produtos';

  constructor(private http: HttpClient) {}

  // Permite cadastrar um produto previamente [cite: 24]
  salvarProduto(novoProduto: {
    codigo: string;
    descricao: string;
    saldo: number;
  }): Observable<any> {
    return this.http.post(this.apiUrl, novoProduto);
  }

  // Lista produtos para utilização em notas fiscais [cite: 25]
  listarProdutos(): Observable<any[]> {
    return this.http.get<any[]>(this.apiUrl);
  }

  // Método de reposição (Soma ao estoque)
  adicionarEstoque(id: number, quantidade: number): Observable<any> {
    return this.http.patch(`${this.apiUrl}/${id}/adicionar`, { quantidade });
  }
}
