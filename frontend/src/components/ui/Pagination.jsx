import { ChevronLeft, ChevronRight } from 'lucide-react';

export function Pagination({ meta, onPageChange }) {
    if (!meta || (meta.last_page <= 1 && meta.total_pages <= 1)) return null;

    const page = meta.current_page || meta.page || 1;
    const total_pages = meta.last_page || meta.total_pages || 1;
    const total_count = meta.total_data || meta.total_count || 0;
    const per_page = meta.per_page || 20;

    const start = (page - 1) * per_page + 1;
    const end = Math.min(page * per_page, total_count);

    const pages = [];
    if (total_pages <= 7) {
        for (let i = 1; i <= total_pages; i++) {
            pages.push(i);
        }
    } else {
        if (page > 3) pages.push(1);
        if (page > 4) pages.push('...');

        let startPage = Math.max(1, page - 2);
        let endPage = Math.min(total_pages, page + 2);

        if (page <= 3) {
            startPage = 1;
            endPage = 5;
        }
        if (page >= total_pages - 2) {
            startPage = total_pages - 4;
            endPage = total_pages;
        }

        for (let i = startPage; i <= endPage; i++) {
            pages.push(i);
        }

        if (page < total_pages - 3) pages.push('...');
        if (page < total_pages - 2) pages.push(total_pages);
    }

    const buttonBase = `
        inline-flex items-center justify-center
        h-9 min-w-[36px] px-3
        text-sm font-medium
        border border-base-300
        transition-colors duration-150
        focus:outline-none focus-visible:ring-2 focus-visible:ring-primary/30
    `;

    return (
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4 pt-4 border-t border-base-200">
            <p className="text-sm text-base-content/60">
                Showing <span className="font-medium text-base-content">{start}</span> to <span className="font-medium text-base-content">{end}</span> of{' '}
                <span className="font-medium text-base-content">{total_count}</span> results
            </p>

            <div className="flex items-center">
                <button
                    className={`${buttonBase} rounded-l-md bg-base-100 hover:bg-base-200 disabled:opacity-50 disabled:cursor-not-allowed`}
                    disabled={page === 1}
                    onClick={() => onPageChange(page - 1)}
                >
                    <ChevronLeft size={16} />
                </button>

                {pages.map((p, idx) => (
                    <button
                        key={idx}
                        className={`
                            ${buttonBase}
                            -ml-px
                            ${p === page ? 'bg-primary text-primary-content border-primary z-10' : 'bg-base-100 hover:bg-base-200'}
                            ${typeof p !== 'number' ? 'cursor-default' : ''}
                        `}
                        onClick={() => typeof p === 'number' && onPageChange(p)}
                        disabled={typeof p !== 'number'}
                    >
                        {p}
                    </button>
                ))}

                <button
                    className={`${buttonBase} -ml-px rounded-r-md bg-base-100 hover:bg-base-200 disabled:opacity-50 disabled:cursor-not-allowed`}
                    disabled={page === total_pages}
                    onClick={() => onPageChange(page + 1)}
                >
                    <ChevronRight size={16} />
                </button>
            </div>
        </div>
    );
}
