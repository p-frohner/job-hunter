const HIDDEN_URLS_KEY = "jh_hidden_job_urls";
const MAX_HIDDEN_URLS = 500;

export const getHiddenUrls = (): Set<string> => {
	try {
		return new Set(JSON.parse(localStorage.getItem(HIDDEN_URLS_KEY) ?? "[]"));
	} catch {
		return new Set();
	}
};

export const addHiddenUrl = (url: string) => {
	const set = getHiddenUrls();
	set.add(url);
	const entries = [...set];
	localStorage.setItem(
		HIDDEN_URLS_KEY,
		JSON.stringify(entries.length > MAX_HIDDEN_URLS ? entries.slice(-MAX_HIDDEN_URLS) : entries),
	);
};
