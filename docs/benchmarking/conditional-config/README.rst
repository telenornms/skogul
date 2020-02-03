Benchmarking Transformer Configuration
======================================

We were curious about how much of a performance hit running
*all* transformers regardless of if they do anything or not
would be.

Initially, we ran a handler with 0 transformers, after that
we ran with 11 transformers and finally we ran with a
conditional transformer for each of the 11 transformers
where none of them applies (the case is false for all).

This is to benchmark how a conditional transformer
theoretically could increase throughput rather than
going throguh every single transformer.

The following configurations were experimented with:

- None: No transformers (Control) (Result: ``none.json``)
- All: All 11 transformers (Result: ``all.json``)
- Conditional: Conditional checks for 11 transformers
  (Result: ``conditional.json``)

The benchmarks ran on an Intel Core i9-9980HK CPU @ 2.40GHz.

The experiments were run on a built skogul binary with
the following command: ``./skogul -f benchmark-config.json``

For the different experiments, the "handler" in
``benchmark-config.json`` test receiver was changed between the
available transformers (``none``, ``all`` or ``conditional``).

All experiments ran for 10 seconds, each producing 10 results.

The following table shows the results.

+---------------+--------------+----------------+--------+-------------------+
| Configuration | Total values | Average rate/s | Ops/ms | Compared to ctrl  |
+===============+==============+================+========+===================+
| None          |    140306230 |       14030623 |  14030 |                \- |
+---------------+--------------+----------------+--------+-------------------+
| All           |     70478220 |        7047822 |   7048 |            50.23% |
+---------------+--------------+----------------+--------+-------------------+
| Conditional   |    136403880 |       13640388 |  13640 |            97.22% |
+---------------+--------------+----------------+--------+-------------------+

None is used as the control and does 14030 operations per millisecond
when the data just passes through. By using a conditional transformer
which skips all transformers we achieve 97% of this throughput, at
13640 operations per millisecond. By executing all transformers
regardless of if they're needed or not (and they're all no-ops)
we achieve 7048 operations per millisecond, which is 50% of the
original throughput.

**Conclusion**: Use the switch transformer to conditionally apply
transformers when performance is paramount.
