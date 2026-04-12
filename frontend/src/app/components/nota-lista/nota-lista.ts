import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FaturamentoService } from '../../services/faturamento';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-nota-lista',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './nota-lista.html',
  styleUrls: ['./nota-lista.scss'],
})
export class NotaListaComponent implements OnInit {
  notas: any[] = [];
  produtosEstoque: any[] = [];
  itensNovaNota: any[] = [];

  notaAbertaId: number | null = null;
  carregandoId: number | null = null;

  constructor(
    private faturamentoService: FaturamentoService,
    private cdr: ChangeDetectorRef,
  ) {}

  ngOnInit(): void {
    this.carregarNotas();
    this.carregarProdutosParaSelecao();
  }

  carregarNotas() {
    this.faturamentoService.listarNotas().subscribe({
      next: (dados: any) => {
        this.notas = dados;
        this.cdr.detectChanges();
      },
      error: (err: any) => console.error('Erro ao buscar notas:', err),
    });
  }

  carregarProdutosParaSelecao() {
    this.faturamentoService.listarProdutosEstoque().subscribe({
      next: (produtos: any) => {
        this.produtosEstoque = produtos;
        this.cdr.detectChanges();
      },
      error: (err) => console.error('Erro ao carregar produtos do estoque:', err),
    });
  }

  adicionarItemLista(
    idProd: string,
    qtd: string,
    prodSel: HTMLSelectElement,
    qtdItem: HTMLInputElement,
  ) {
    const produtoId = parseInt(idProd);
    const quantidade = parseInt(qtd);

    if (isNaN(produtoId) || quantidade <= 0) {
      alert('Selecione um produto e uma quantidade válida!');
      return;
    }

    // Busca o produto na lista carregada (propriedades minúsculas do estoque)
    const produtoOriginal = this.produtosEstoque.find((p) => p.id === produtoId);

    if (produtoOriginal) {
      this.itensNovaNota.push({
        ProdutoId: produtoId,
        Quantidade: quantidade,
        descricao_temporaria: produtoOriginal.descricao,
      });
      prodSel.value = '';
      qtdItem.value = '1';
      this.cdr.detectChanges();
    } else {
      alert('Produto não encontrado na lista local.');
    }

    this.cdr.detectChanges();
  }

  removerItemLista(index: number) {
    this.itensNovaNota.splice(index, 1);
    this.cdr.detectChanges();
  }

  criarNota(numero: string) {
    const num = parseInt(numero);
    if (isNaN(num) || num <= 0) {
      alert('O Número Sequencial é obrigatório e deve ser maior que zero!');
      return;
    }

    if (this.itensNovaNota.length === 0) {
      alert('Adicione pelo menos um produto à nota antes de finalizar!');
      return;
    }

    // Objeto montado para casar com a Struct Go (PascalCase)
    const novaNota = {
      NumeroSequencial: num,
      Status: 'Aberta',
      Itens: this.itensNovaNota.map((i) => ({
        ProdutoId: i.ProdutoId,
        Quantidade: i.Quantidade,
      })),
    };

    console.log('ENVIANDO PARA O GO:', novaNota);

    this.faturamentoService.abrirNota(novaNota).subscribe({
      next: () => {
        alert('Nota Fiscal nº ' + num + ' criada com sucesso!');
        this.itensNovaNota = [];
        this.carregarNotas();
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error('Erro retornado do Go:', err);
        alert('Falha no Backend: ' + (err.error?.error || 'Erro desconhecido'));
      },
    });
  }

  toggleDetalhes(id: number) {
    this.notaAbertaId = this.notaAbertaId === id ? null : id;
    this.cdr.detectChanges();
  }

  // AJUSTE: Recebe nota.Id (PascalCase)
  imprimirNota(id: number) {
    if (!id) return;
    this.carregandoId = id;
    this.cdr.detectChanges();

    this.faturamentoService.fecharNota(id).subscribe({
      next: (res: any) => {
        alert('Sucesso: Nota Fechada e Estoque Atualizado!');
        this.carregandoId = null;
        this.carregarNotas();
      },
      error: (err: any) => {
        this.carregandoId = null;
        const mensagemErro = err.error?.error || 'Erro na integração com o estoque.';
        alert('Falha na Operação: ' + mensagemErro);
        this.cdr.detectChanges();
      },
    });
  }
}
