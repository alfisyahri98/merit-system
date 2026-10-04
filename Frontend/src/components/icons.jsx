// Ikon garis sederhana (24x24, stroke). Ukuran diatur lewat className.
function Svg({ children, className = "size-5" }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" className={className} aria-hidden="true">
      {children}
    </svg>
  );
}

export const IconHome = (p) => <Svg {...p}><path d="M3 10.5 12 3l9 7.5V20a1 1 0 0 1-1 1h-5v-6h-6v6H4a1 1 0 0 1-1-1z" /></Svg>;
export const IconUsers = (p) => <Svg {...p}><circle cx="9" cy="8" r="3.5" /><path d="M2.5 20c0-3.6 2.9-6 6.5-6s6.5 2.4 6.5 6" /><path d="M16 4.5a3.5 3.5 0 0 1 0 7M18 14c2.2.6 3.5 2.8 3.5 6" /></Svg>;
export const IconUserCog = (p) => <Svg {...p}><circle cx="10" cy="8" r="3.5" /><path d="M3 20c0-3.6 3-6 7-6 1 0 1.9.1 2.7.4" /><circle cx="18" cy="17" r="2.5" /><path d="M18 13v1.5M18 19.5V21M14 17h1.5M20.5 17H22" /></Svg>;
export const IconKey = (p) => <Svg {...p}><circle cx="8" cy="15" r="4" /><path d="m11 12 9-9M17 6l3 3M14 9l2 2" /></Svg>;
export const IconLogout = (p) => <Svg {...p}><path d="M15 4h3a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2h-3M10 17l-5-5 5-5M5 12h11" /></Svg>;
export const IconMenu = (p) => <Svg {...p}><path d="M4 6h16M4 12h16M4 18h16" /></Svg>;
export const IconX = (p) => <Svg {...p}><path d="M6 6l12 12M18 6 6 18" /></Svg>;
export const IconSearch = (p) => <Svg {...p}><circle cx="11" cy="11" r="7" /><path d="m20 20-3.5-3.5" /></Svg>;
export const IconPlus = (p) => <Svg {...p}><path d="M12 5v14M5 12h14" /></Svg>;
export const IconPencil = (p) => <Svg {...p}><path d="M4 20h4L19 9l-4-4L4 16z" /></Svg>;
export const IconTrash = (p) => <Svg {...p}><path d="M4 7h16M10 11v6M14 11v6M6 7l1 13h10l1-13M9 7V4h6v3" /></Svg>;
export const IconChevronDown = (p) => <Svg {...p}><path d="m6 9 6 6 6-6" /></Svg>;
export const IconArrowLeft = (p) => <Svg {...p}><path d="M19 12H5M11 6l-6 6 6 6" /></Svg>;
export const IconBriefcase = (p) => <Svg {...p}><rect x="3" y="7" width="18" height="13" rx="2" /><path d="M8 7V5a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M3 13h18" /></Svg>;
export const IconCheck = (p) => <Svg {...p}><path d="m5 12 5 5 9-10" /></Svg>;
export const IconAlert = (p) => <Svg {...p}><circle cx="12" cy="12" r="9" /><path d="M12 8v5M12 16h.01" /></Svg>;
export const IconBook = (p) => <Svg {...p}><path d="M4 5a2 2 0 0 1 2-2h13v16H6a2 2 0 0 0-2 2zM4 19V5M9 7h6" /></Svg>;
export const IconDownload = (p) => <Svg {...p}><path d="M12 4v11M7 10l5 5 5-5M5 20h14" /></Svg>;
export const IconFilter = (p) => <Svg {...p}><path d="M4 5h16l-6 7.5V19l-4 2v-8.5z" /></Svg>;
export const IconArrowUp = (p) => <Svg {...p}><path d="M12 19V5M6 11l6-6 6 6" /></Svg>;
export const IconArrowDown = (p) => <Svg {...p}><path d="M12 5v14M6 13l6 6 6-6" /></Svg>;
export const IconSort = (p) => <Svg {...p}><path d="M8 4v16M4 8l4-4 4 4M16 20V4M12 16l4 4 4-4" /></Svg>;
