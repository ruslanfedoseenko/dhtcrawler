package Models

type PaginatedTorrentsResponse struct {
	Page         uint64
	PageCount    uint64
	ItemsCount   uint64
	ItemsPerPage uint64
	Torrents     []Torrent
}
