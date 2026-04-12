import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class FaturamentoService {
  // Porta 8082: Microsserviço de Faturamento
  private apiUrl = 'http://localhost:8082/faturamento';

  // Porta 8081: Microsserviço de Estoque
  private estoqueUrl = 'http://localhost:8081/produtos';

  constructor(private http: HttpClient) {}

  /**
   * METODOS DE FATURAMENTO (Porta 8082)
   */

  // Busca todas as notas para listar na tela
  listarNotas(): Observable<any[]> {
    return this.http.get<any[]>(this.apiUrl);
  }

  // Cria a nota com os itens selecionados [Requisito: Cadastro de Notas]
  abrirNota(nota: any): Observable<any> {
    return this.http.post(this.apiUrl, nota);
  }

  // Finaliza a nota e dispara a integração de baixa [Requisito 3.35]
  fecharNota(id: number): Observable<any> {
    return this.http.patch(`${this.apiUrl}/${id}/fechar`, {});
  }

  /**
   * METODOS DE INTEGRAÇÃO/APOIO (Porta 8081)
   */

  // Busca produtos do estoque para preencher o <select> no formulário de faturamento
  listarProdutosEstoque(): Observable<any[]> {
    return this.http.get<any[]>(this.estoqueUrl);
  }
}
