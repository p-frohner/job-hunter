import { useCallback, useRef, useState } from "react";
import type { Job, SourceResult } from "../api";
import { addHiddenUrl, getHiddenUrls } from "../utils/hiddenJobs";

export function useSearchStream() {
	const [jobs, setJobs] = useState<Job[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const esRef = useRef<EventSource | null>(null);

	const search = useCallback((query: string, location?: string) => {
		esRef.current?.close();
		setJobs([]);
		setError(null);
		setLoading(true);

		const params = new URLSearchParams({ query });
		if (location) {
			params.set("location", location);
		}

		const es = new EventSource(`${import.meta.env.VITE_API_URL}/api/search/stream?${params}`);
		esRef.current = es;

		es.onmessage = (e) => {
			let data: SourceResult & { done?: boolean };
			try {
				data = JSON.parse(e.data);
			} catch {
				setError("Received malformed data from server");
				setLoading(false);
				es.close();
				return;
			}

			if (data.done) {
				setLoading(false);
				es.close();
				return;
			}

			const hiddenUrls = getHiddenUrls();
			setJobs((prev) => [...prev, ...(data.jobs ?? []).filter((j) => !hiddenUrls.has(j.url))]);
		};

		es.onerror = () => {
			setError("Stream error — check the server logs");
			setLoading(false);
			es.close();
		};
	}, []);

	const removeJob = useCallback((id: string, source: string | undefined, url: string) => {
		addHiddenUrl(url);
		setJobs((prev) => prev.filter((j) => !(j.id === id && j.source === source)));
	}, []);

	return { search, jobs, loading, error, removeJob };
}
