class SearchStore {
	focusSeq = $state(0);

	requestFocus() {
		this.focusSeq++;
	}
}

export const searchStore = new SearchStore();
