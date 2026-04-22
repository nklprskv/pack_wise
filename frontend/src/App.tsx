import { useEffect, useMemo, useState } from "react";
import { Calculator, Plus, RefreshCw, ShieldCheck, Trash2 } from "lucide-react";

import { apiFetch } from "./lib/api";

type CalculationPack = {
  size: number;
  quantity: number;
};

type CalculationResult = {
  requestedItems: number;
  totalItems: number;
  packs: CalculationPack[];
};

const defaultPackSizes = [250, 500, 1000, 2000, 5000];
const challengePackSizes = [23, 31, 53];

export default function App() {
  const [packSizes, setPackSizes] = useState<number[]>([]);
  const [packSizeDraft, setPackSizeDraft] = useState("");
  const [items, setItems] = useState("501");
  const [result, setResult] = useState<CalculationResult | null>(null);
  const [isBusy, setIsBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const sortedPackSizes = useMemo(
    () => [...packSizes].sort((left, right) => left - right),
    [packSizes],
  );

  useEffect(() => {
    void loadPackSizes();
  }, []);

  async function loadPackSizes() {
    setIsBusy(true);
    setError(null);

    try {
      const response = await apiFetch("/api/v1/packs");
      if (!response.ok) {
        throw new Error("Failed to load pack sizes");
      }

      const data = (await response.json()) as number[];
      setPackSizes(data);
    } catch (requestError) {
      setError(getErrorMessage(requestError));
    } finally {
      setIsBusy(false);
    }
  }

  async function replacePackSizes(nextPackSizes: number[]) {
    setIsBusy(true);
    setError(null);

    try {
      const response = await apiFetch("/api/v1/packs", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ sizes: nextPackSizes }),
      });

      if (!response.ok) {
        throw new Error(await readErrorDetail(response, "Failed to update pack sizes"));
      }

      setPackSizes(nextPackSizes);
    } catch (requestError) {
      setError(getErrorMessage(requestError));
    } finally {
      setIsBusy(false);
    }
  }

  async function deletePackSize(size: number) {
    setIsBusy(true);
    setError(null);

    try {
      const response = await apiFetch(`/api/v1/packs/${size}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error(await readErrorDetail(response, "Failed to delete pack size"));
      }

      setPackSizes((currentPackSizes) => currentPackSizes.filter((currentSize) => currentSize !== size));
    } catch (requestError) {
      setError(getErrorMessage(requestError));
    } finally {
      setIsBusy(false);
    }
  }

  async function handleCalculate() {
    setIsBusy(true);
    setError(null);
    setResult(null);

    try {
      const response = await apiFetch("/api/v1/calculate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ items: Number(items) }),
      });

      if (!response.ok) {
        throw new Error(await readErrorDetail(response, "Failed to calculate packs"));
      }

      const data = (await response.json()) as CalculationResult;
      setResult(data);
    } catch (requestError) {
      setError(getErrorMessage(requestError));
    } finally {
      setIsBusy(false);
    }
  }

  async function handleAddPackSize() {
    const nextSize = Number(packSizeDraft);
    if (!Number.isFinite(nextSize) || nextSize <= 0) {
      setError("Pack size must be a positive number.");
      return;
    }

    if (packSizes.includes(nextSize)) {
      setError("Pack size already exists.");
      return;
    }

    const nextPackSizes = [...packSizes, nextSize].sort((left, right) => left - right);
    setPackSizeDraft("");
    await replacePackSizes(nextPackSizes);
  }

  async function handleResetDefaultSizes() {
    await replacePackSizes(defaultPackSizes);
  }

  async function handleChallengeCase() {
    setItems("500000");
    await replacePackSizes(challengePackSizes);
  }

  return (
    <main className="min-h-screen px-4 py-8 text-ink sm:px-6 lg:px-8">
      <div className="mx-auto max-w-7xl">
        <header className="mb-8 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
          <div className="space-y-3">
            <div className="inline-flex items-center gap-2 rounded-full border border-black/5 bg-white/80 px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-steel shadow-panel">
              <ShieldCheck className="h-4 w-4 text-pine" />
              Packaging Optimizer
            </div>
            <div className="space-y-2">
              <h1 className="text-4xl font-semibold tracking-tight sm:text-5xl">PackWise dashboard</h1>
              <p className="max-w-2xl text-sm leading-6 text-steel sm:text-base">
                Manage active pack sizes and run fulfillment calculations against the current configuration.
              </p>
            </div>
          </div>

          <div className="flex flex-wrap gap-3">
            <button
              className="inline-flex items-center gap-2 rounded-full border border-ink/10 bg-white px-4 py-2 text-sm font-medium text-ink shadow-panel transition hover:border-ink/20 hover:bg-white/80"
              onClick={() => void handleResetDefaultSizes()}
              type="button"
            >
              <RefreshCw className="h-4 w-4" />
              Reset Defaults
            </button>
            <button
              className="inline-flex items-center gap-2 rounded-full bg-ink px-4 py-2 text-sm font-medium text-white transition hover:bg-black"
              onClick={() => void handleChallengeCase()}
              type="button"
            >
              <Calculator className="h-4 w-4" />
              Challenge Case
            </button>
          </div>
        </header>

        {error ? (
          <div className="mb-6 rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        ) : null}

        <section className="grid gap-6 lg:grid-cols-[1fr,1.2fr]">
          <article className="rounded-[28px] border border-black/5 bg-white/90 p-6 shadow-panel backdrop-blur">
            <div className="mb-6 flex items-center justify-between">
              <div>
                <p className="text-xs font-semibold uppercase tracking-[0.24em] text-steel">Pack Sizes</p>
                <h2 className="mt-2 text-2xl font-semibold">Current configuration</h2>
              </div>
              <span className="rounded-full bg-ink px-3 py-1 text-xs font-semibold text-white">
                {sortedPackSizes.length} active
              </span>
            </div>

            <div className="mb-5 flex gap-3">
              <input
                className="min-w-0 flex-1 rounded-2xl border border-black/10 bg-sand px-4 py-3 text-sm outline-none transition focus:border-ink/40"
                inputMode="numeric"
                onChange={(event) => setPackSizeDraft(event.target.value)}
                placeholder="Add pack size"
                value={packSizeDraft}
              />
              <button
                className="inline-flex items-center gap-2 rounded-2xl bg-flare px-4 py-3 text-sm font-medium text-white transition hover:bg-amber-700"
                disabled={isBusy}
                onClick={() => void handleAddPackSize()}
                type="button"
              >
                <Plus className="h-4 w-4" />
                Add
              </button>
            </div>

            <div className="space-y-3">
              {sortedPackSizes.length === 0 ? (
                <div className="rounded-2xl border border-dashed border-black/10 bg-sand px-4 py-6 text-sm text-steel">
                  No pack sizes configured.
                </div>
              ) : (
                sortedPackSizes.map((size) => (
                  <div
                    className="flex items-center justify-between rounded-2xl border border-black/5 bg-sand px-4 py-3"
                    key={size}
                  >
                    <div>
                      <p className="text-xs uppercase tracking-[0.24em] text-steel">Pack size</p>
                      <p className="mt-1 text-lg font-semibold">{size}</p>
                    </div>
                    <button
                      className="inline-flex items-center gap-2 rounded-full border border-red-200 px-3 py-2 text-sm font-medium text-red-600 transition hover:bg-red-50"
                      disabled={isBusy}
                      onClick={() => void deletePackSize(size)}
                      type="button"
                    >
                      <Trash2 className="h-4 w-4" />
                      Delete
                    </button>
                  </div>
                ))
              )}
            </div>
          </article>

          <article className="rounded-[28px] border border-black/5 bg-white/90 p-6 shadow-panel backdrop-blur">
            <div className="mb-6">
              <p className="text-xs font-semibold uppercase tracking-[0.24em] text-steel">Calculator</p>
              <h2 className="mt-2 text-2xl font-semibold">Shipment breakdown</h2>
            </div>

            <div className="mb-6 flex flex-col gap-3 sm:flex-row">
              <input
                className="min-w-0 flex-1 rounded-2xl border border-black/10 bg-sand px-4 py-3 text-sm outline-none transition focus:border-ink/40"
                inputMode="numeric"
                onChange={(event) => setItems(event.target.value)}
                placeholder="Requested items"
                value={items}
              />
              <button
                className="inline-flex items-center justify-center gap-2 rounded-2xl bg-ink px-5 py-3 text-sm font-medium text-white transition hover:bg-black"
                disabled={isBusy}
                onClick={() => void handleCalculate()}
                type="button"
              >
                <Calculator className="h-4 w-4" />
                Calculate
              </button>
            </div>

            {result ? (
              <div className="space-y-5">
                <div className="grid gap-3 sm:grid-cols-2">
                  <MetricCard label="Requested items" value={result.requestedItems.toLocaleString()} />
                  <MetricCard label="Shipped items" value={result.totalItems.toLocaleString()} />
                </div>

                <div className="overflow-hidden rounded-3xl border border-black/5">
                  <table className="min-w-full divide-y divide-black/5 text-left text-sm">
                    <thead className="bg-sand text-steel">
                      <tr>
                        <th className="px-4 py-3 font-medium">Pack size</th>
                        <th className="px-4 py-3 font-medium">Quantity</th>
                        <th className="px-4 py-3 font-medium">Items shipped</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-black/5 bg-white">
                      {result.packs.map((pack) => (
                        <tr key={`${pack.size}-${pack.quantity}`}>
                          <td className="px-4 py-3 font-medium">{pack.size}</td>
                          <td className="px-4 py-3">{pack.quantity}</td>
                          <td className="px-4 py-3">{(pack.size * pack.quantity).toLocaleString()}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            ) : (
              <div className="rounded-3xl border border-dashed border-black/10 bg-sand px-6 py-12 text-center text-sm text-steel">
                Run a calculation to see the pack breakdown.
              </div>
            )}
          </article>
        </section>
      </div>
    </main>
  );
}

function MetricCard({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-3xl border border-black/5 bg-sand px-5 py-4">
      <p className="text-xs uppercase tracking-[0.24em] text-steel">{label}</p>
      <p className="mt-2 text-2xl font-semibold">{value}</p>
    </div>
  );
}

async function readErrorDetail(response: Response, fallbackMessage: string) {
  const contentType = response.headers.get("Content-Type") ?? "";
  if (!contentType.includes("application/json")) {
    return fallbackMessage;
  }

  const data = (await response.json()) as { detail?: string };
  return data.detail ?? fallbackMessage;
}

function getErrorMessage(error: unknown) {
  if (error instanceof Error) {
    return error.message;
  }

  return "Unexpected error";
}
