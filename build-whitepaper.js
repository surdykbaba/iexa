const fs = require("fs");
const {
  Document, Packer, Paragraph, TextRun, Table, TableRow, TableCell,
  AlignmentType, LevelFormat, HeadingLevel, BorderStyle, WidthType, ShadingType,
  TableOfContents, PageNumber, PageBreak, Header, Footer, VerticalAlign,
} = require("docx");

const INDIGO = "0F2C4D", TEAL = "0F9D8F", GOLD = "B8860B", INK = "101828", MUTED = "5B6472";
const CW = 9360; // content width (US Letter, 1" margins)

const border = { style: BorderStyle.SINGLE, size: 1, color: "DDE2EA" };
const borders = { top: border, bottom: border, left: border, right: border,
  insideHorizontal: border, insideVertical: border };

// ---- helpers ----
function h1(text) {
  return new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun(text)] });
}
function h2(text) {
  return new Paragraph({ heading: HeadingLevel.HEADING_2, children: [new TextRun(text)] });
}
function p(text, opts = {}) {
  return new Paragraph({
    spacing: { after: 160, line: 276 },
    children: [new TextRun({ text, ...opts })],
  });
}
function rich(runs) {
  return new Paragraph({ spacing: { after: 160, line: 276 }, children: runs });
}
function bullet(text, bold) {
  return new Paragraph({
    numbering: { reference: "bullets", level: 0 },
    spacing: { after: 80, line: 276 },
    children: bold
      ? [new TextRun({ text: bold, bold: true }), new TextRun({ text: " — " + text })]
      : [new TextRun(text)],
  });
}
function cell(text, { headerRow = false, w, bold = false, fill } = {}) {
  return new TableCell({
    borders,
    width: { size: w, type: WidthType.DXA },
    margins: { top: 80, bottom: 80, left: 120, right: 120 },
    shading: fill ? { fill, type: ShadingType.CLEAR } : (headerRow ? { fill: INDIGO, type: ShadingType.CLEAR } : undefined),
    verticalAlign: VerticalAlign.CENTER,
    children: [new Paragraph({
      spacing: { after: 0, line: 264 },
      children: [new TextRun({ text, bold: headerRow || bold, color: headerRow ? "FFFFFF" : INK, size: 20 })],
    })],
  });
}
function table(rows, widths) {
  return new Table({
    width: { size: CW, type: WidthType.DXA },
    columnWidths: widths,
    rows: rows.map((r, i) =>
      new TableRow({
        tableHeader: i === 0,
        children: r.map((c, j) => cell(c, { headerRow: i === 0, w: widths[j], fill: i > 0 && j === 0 ? "F2F5F9" : undefined })),
      })
    ),
  });
}
function callout(title, text) {
  return new Table({
    width: { size: CW, type: WidthType.DXA },
    columnWidths: [CW],
    rows: [new TableRow({ children: [new TableCell({
      borders: { left: { style: BorderStyle.SINGLE, size: 24, color: TEAL },
        top: { style: BorderStyle.NONE }, bottom: { style: BorderStyle.NONE }, right: { style: BorderStyle.NONE } },
      width: { size: CW, type: WidthType.DXA },
      shading: { fill: "F0FAF8", type: ShadingType.CLEAR },
      margins: { top: 140, bottom: 140, left: 200, right: 160 },
      children: [
        new Paragraph({ spacing: { after: 60 }, children: [new TextRun({ text: title, bold: true, color: INDIGO, size: 22 })] }),
        new Paragraph({ spacing: { after: 0, line: 276 }, children: [new TextRun({ text, size: 21 })] }),
      ],
    })] })],
  });
}
function spacer(after = 120) { return new Paragraph({ spacing: { after }, children: [] }); }

// ---- document ----
const doc = new Document({
  creator: "IXEA",
  title: "IXEA Whitepaper",
  description: "Interoperability and Exchange Alliance — Concept & Founding Framework",
  styles: {
    default: { document: { run: { font: "Arial", size: 21, color: INK } } },
    paragraphStyles: [
      { id: "Heading1", name: "Heading 1", basedOn: "Normal", next: "Normal", quickFormat: true,
        run: { size: 30, bold: true, color: INDIGO, font: "Arial" },
        paragraph: { spacing: { before: 320, after: 160 }, outlineLevel: 0,
          border: { bottom: { style: BorderStyle.SINGLE, size: 6, color: TEAL, space: 6 } } } },
      { id: "Heading2", name: "Heading 2", basedOn: "Normal", next: "Normal", quickFormat: true,
        run: { size: 24, bold: true, color: INK, font: "Arial" },
        paragraph: { spacing: { before: 220, after: 120 }, outlineLevel: 1 } },
    ],
  },
  numbering: {
    config: [
      { reference: "bullets", levels: [{ level: 0, format: LevelFormat.BULLET, text: "•", alignment: AlignmentType.LEFT,
        style: { paragraph: { indent: { left: 540, hanging: 280 } } } }] },
      { reference: "nums", levels: [{ level: 0, format: LevelFormat.DECIMAL, text: "%1.", alignment: AlignmentType.LEFT,
        style: { paragraph: { indent: { left: 540, hanging: 280 } } } }] },
    ],
  },
  sections: [
    // ---------- COVER ----------
    {
      properties: { page: { size: { width: 12240, height: 15840 }, margin: { top: 1440, right: 1440, bottom: 1440, left: 1440 } } },
      children: [
        spacer(1600),
        new Paragraph({ alignment: AlignmentType.LEFT, spacing: { after: 80 },
          children: [new TextRun({ text: "IXEA", bold: true, size: 84, color: INDIGO })] }),
        new Paragraph({ spacing: { after: 400 },
          children: [new TextRun({ text: "Interoperability and Exchange Alliance", size: 32, color: TEAL, bold: true })] }),
        new Paragraph({ spacing: { after: 120 },
          children: [new TextRun({ text: "Together towards a federated and sovereign", size: 40, bold: true, color: INK })] }),
        new Paragraph({ spacing: { after: 480 },
          children: [new TextRun({ text: "data infrastructure for Africa.", size: 40, bold: true, color: INK })] }),
        new Paragraph({ spacing: { after: 120 },
          children: [new TextRun({ text: "Concept & Founding Framework", size: 26, color: MUTED })] }),
        new Paragraph({ spacing: { after: 80 },
          children: [new TextRun({ text: "Version 0.1  ·  Working Draft for Community Review", size: 22, color: MUTED, italics: true })] }),
        spacer(2400),
        new Paragraph({ border: { top: { style: BorderStyle.SINGLE, size: 6, color: "DDE2EA", space: 8 } },
          spacing: { before: 200, after: 60 },
          children: [new TextRun({ text: "An open, community-governed standard. Free and open-source under Apache-2.0.", size: 20, color: MUTED })] }),
        new Paragraph({ children: [new TextRun({ text: "Pan-African · Aspiring Digital Public Good", size: 20, color: MUTED })] }),
      ],
    },
    // ---------- BODY ----------
    {
      properties: { page: { size: { width: 12240, height: 15840 }, margin: { top: 1440, right: 1440, bottom: 1440, left: 1440 } } },
      footers: { default: new Footer({ children: [new Paragraph({
        alignment: AlignmentType.CENTER,
        border: { top: { style: BorderStyle.SINGLE, size: 4, color: "DDE2EA", space: 6 } },
        children: [new TextRun({ text: "IXEA · Interoperability and Exchange Alliance · v0.1   |   Page ", size: 16, color: MUTED }),
          new TextRun({ children: [PageNumber.CURRENT], size: 16, color: MUTED })] })] }) },
      children: [
        new Paragraph({ spacing: { after: 200 }, children: [new TextRun({ text: "Table of Contents", bold: true, size: 28, color: INDIGO })] }),
        new TableOfContents("Table of Contents", { hyperlink: true, headingStyleRange: "1-2" }),
        new Paragraph({ children: [new PageBreak()] }),

        // 1. Executive Summary
        h1("1. Executive Summary"),
        p("Africa's digital economy is being built on fragmented, disconnected systems. Every bank, fintech, government agency and ERP integrates point-to-point with every other — a costly, brittle web of bespoke connections that locks data in silos, raises the cost of doing business, and leaves the continent's most valuable asset, its data, under foreign control."),
        rich([
          new TextRun("IXEA — the "),
          new TextRun({ text: "Interoperability and Exchange Alliance", bold: true }),
          new TextRun(" — is an open, community-governed framework for secure, federated data exchange across the continent. It is not a product and not owned by any single company. It is a neutral standard, a trust framework, and a community — modelled on the world's most successful interoperability initiatives: Estonia's X-Road, the EU's Gaia-X and Peppol, and the MOSIP identity platform."),
        ]),
        p("This document sets out the concept, the architecture, the governance model, and the roadmap. It is a working draft, published openly to invite the developers, providers, regulators and young technologists who will build it."),
        callout("In one sentence", "IXEA lets any organisation connect once and securely exchange invoices, identity and trusted data with everyone else on the network — without rebuilding integrations, and without surrendering data sovereignty."),

        // 2. The Problem
        h1("2. The Problem"),
        p("Interoperability is the missing layer of Africa's digital infrastructure. The symptoms are everywhere:"),
        bullet("the same KYC document is collected, re-keyed and re-verified by every institution a citizen touches.", "Duplicated effort"),
        bullet("invoices, identity attributes and records are trapped in incompatible formats and proprietary platforms.", "Data silos"),
        bullet("each new partner integration is a months-long, point-to-point engineering project.", "Integration tax"),
        bullet("critical national data flows through infrastructure owned and operated outside the continent.", "Lost sovereignty"),
        bullet("without a shared trust layer, fraud thrives and cross-border commerce stalls.", "Weak trust"),
        p("These are not problems any one company or country can solve alone. They are coordination problems — and coordination problems are solved by shared, open standards governed by a neutral body."),

        // 3. Vision & Mission
        h1("3. Vision & Mission"),
        h2("Vision"),
        p("A federated, trusted and sovereign data ecosystem that connects every African institution — public and private, across all 54 countries — so that data can move securely to where it creates value, while remaining under African control."),
        h2("Mission"),
        p("To create and steward the open standards, trust framework and reference software that enable interoperable data exchange across Africa, and to grow the community of technologists who build and run it."),
        h2("Principles"),
        bullet("the standard belongs to the community, not to any company or government.", "Neutral"),
        bullet("specifications, code and governance are public and contributable by anyone.", "Open"),
        bullet("data stays at its source; the network moves messages, not a central data lake.", "Sovereign"),
        bullet("every participant is verifiable; every exchange is signed and auditable.", "Trusted"),
        bullet("designed to interoperate with global standards, never to reinvent them.", "Interoperable"),

        // 4. What is IXEA
        h1("4. What IXEA Is"),
        p("IXEA operates on three distinct but connected layers. Keeping them separate is what makes the initiative durable."),
        table([
          ["Layer", "What it is", "Comparable to"],
          ["The Association", "The neutral non-profit that sets rules, admits members and accredits providers.", "OpenPeppol AISBL, NIIS (X-Road)"],
          ["The Standard", "The open specifications: how identity, invoices and data are described, signed and transported.", "Peppol BIS, X-Road protocol"],
          ["Reference Software", "Open-source nodes and a public sandbox that providers run to join the network.", "X-Road core, MOSIP"],
        ], [2200, 4760, 2400]),
        spacer(80),
        callout("Why this separation matters", "The value and defensibility of IXEA is in the standard. Whoever stewards the trusted standard for African data exchange anchors the ecosystem. The Association exists to keep that standard neutral and legitimate; the software exists to make adoption easy."),

        // 5. The Framework — three pillars
        h1("5. The Framework: Three Pillars"),
        h2("Pillar 1 — Open Standards"),
        p("Public, versioned specifications for how data is described, signed and transported, developed through an open RFC process. Anyone can read, use, or propose changes. The standard leads with one use case — e-invoicing — and extends through additional profiles."),
        h2("Pillar 2 — Trust Framework"),
        p("A verification and accreditation framework with a conformance test suite. Every participant is provably who they claim to be; every message is cryptographically signed and auditable. “IXEA-certified” becomes a meaningful badge that guarantees genuine interoperability and security."),
        h2("Pillar 3 — Federation"),
        p("A decentralised network of interoperable nodes and Clearing Houses — no single point of control and no central data store. Trust federation lets independent ecosystems and national hubs connect across borders while each retains its autonomy."),

        // 6. How it works
        h1("6. How It Works"),
        p("Each member runs (or rents) a single secure node — the only component they integrate with. Through IXEA's shared registry and trust services, that one connection reaches every other participant. There are no point-to-point integrations and no central honeypot of data."),
        table([
          ["Step", "Role", "What happens"],
          ["1", "Sender", "A business, bank or agency needs to send data."],
          ["2", "Member Node", "A secure gateway signs, encrypts and routes the message onto IXEA."],
          ["3", "Member Node", "The receiver's node is discovered via the shared registry and trust list."],
          ["4", "Receiver", "The destination acts on it — a paid invoice, a verified identity, an exchanged record."],
        ], [1000, 2400, 5960]),
        spacer(80),
        callout("The Once-Only Principle", "Citizens and businesses should provide the same information to government once — never again. IXEA lets institutions request verified data straight from the source, with consent, instead of asking people to re-submit documents they have already given. Less paperwork, less fraud, more trust."),

        // 7. Data Spaces
        h1("7. Data Spaces"),
        p("IXEA is one framework hosting many data spaces — sectoral ecosystems that share the same trust framework and federation. The strategy is to win one use case with strong regulatory tailwind, then extend."),
        table([
          ["Data space", "What it carries", "Status"],
          ["Invoicing", "Standardised e-invoice exchange between businesses, ERPs and tax authorities — the wedge use case.", "Available · v0.1"],
          ["Identity", "Consent-driven, signed exchange of verified identity attributes and KYC credentials, interoperable with national identity systems.", "Available · v0.1"],
          ["Data (general)", "A general-purpose envelope for any domain — health, logistics, open finance — that adopts the IXEA trust framework.", "On roadmap"],
        ], [2000, 5360, 2000]),

        // 8. Interoperability
        h1("8. Interoperability: Standards & Rails"),
        p("IXEA does not replace the world's standards — it binds them together, so an IXEA invoice can be paid, settled and reported anywhere."),
        h2("Invoice & document standards"),
        p("Peppol BIS · UBL 2.1 · UN/CEFACT CII · EN 16931 · GS1 · national e-invoicing / CTC mandates.", { color: MUTED }),
        h2("Payment & settlement rails"),
        p("EMVCo QR · Request-to-Pay (RTP) · ISO 20022 · national instant-payment switches · PAPSS · Mojaloop · Open Banking APIs.", { color: MUTED }),
        h2("Tax, compliance & incentives"),
        p("VAT / fiscal lottery · continuous transaction controls · digital signatures (XAdES / PKI) · tax-authority clearance · e-receipts.", { color: MUTED }),
        spacer(60),
        callout("Why the VAT lottery matters", "Fiscal receipt lotteries — proven in Portugal, Italy and Greece — reward consumers for requesting a real invoice, turning every buyer into a tax-compliance enforcer. By embedding lottery-eligible receipt identifiers in the invoice profile, any member country can launch a VAT lottery and widen the formal economy from day one."),

        // 9. Governance & Membership
        h1("9. Governance & Membership"),
        p("IXEA is a neutral, not-for-profit association. Governance is shared across its members; no single company owns the standard. Membership is open in three tiers:"),
        table([
          ["Tier", "Who", "What they get"],
          ["Contributor", "Developers, students, researchers", "Free access to specs and sandbox; contribute via the open RFC process; join working groups and hackathons."],
          ["Provider", "Fintechs, vendors, ERPs", "Run certified nodes; earn the IXEA-certified badge; a voting seat in working groups; commercial use under Apache-2.0."],
          ["Authority", "Regulators, public agencies", "Shape policy and national profiles; board representation; accreditation oversight."],
        ], [1700, 3000, 4660]),
        spacer(80),
        h2("Working groups"),
        p("The technical work happens in the open, in working groups: Standards & RFCs, Invoicing, Identity, Transport & Security, and Certification. Each is led by maintainers and shaped by the community."),

        // 10. Community
        h1("10. Community: Bringing Young Minds Together"),
        p("IXEA is as much a movement as a standard. Its purpose is to convene the next generation of African technologists and give them real infrastructure to build. The community is organised around three structures borrowed from the most successful global ecosystems:"),
        bullet("country chapters — starting with Nigeria — that localise the framework, run events and connect members on the ground.", "National Hubs"),
        bullet("free training, certification and mentorship that turns students and young developers into the engineers who run Africa's data infrastructure.", "IXEA Academy"),
        bullet("flagship deployments — a live e-invoicing data space, a cross-border identity pilot — that prove the framework in the real world.", "Lighthouse Projects"),

        // 11. Roadmap
        h1("11. Roadmap"),
        p("A transparent, milestone-driven path — every step open and community-reviewed."),
        table([
          ["Phase", "Milestone", "Focus"],
          ["1 · Now", "Foundation", "Association charter; v0.1 invoicing and identity specs; public sandbox; first member nodes."],
          ["2", "Lighthouse pilots", "Live e-invoicing and cross-border identity deployments with founding members and a national hub."],
          ["3", "Digital Public Good", "Certification with the DPG Alliance — the credibility stamp that unlocks donor and government adoption."],
          ["4", "Pan-African scale", "Trust federation across national hubs — one interoperable data ecosystem, 54 countries."],
        ], [1400, 2600, 5360]),

        // 12. Call to action
        h1("12. Get Involved"),
        p("IXEA will be built by the people who show up. Whether you write code, run infrastructure or set policy, there is a seat for you."),
        bullet("contribute to the specs, run a node in the sandbox, or join a working group.", "Developers"),
        bullet("become a founding member, run a certified node, and help shape the standard.", "Providers & vendors"),
        bullet("co-design national profiles and bring the Once-Only Principle to your citizens.", "Regulators & agencies"),
        bullet("partner on a national Hub, the Academy, or a Lighthouse pilot.", "Universities & funders"),
        spacer(120),
        new Paragraph({ shading: { fill: INDIGO, type: ShadingType.CLEAR }, spacing: { before: 120, after: 120, line: 300 },
          indent: { left: 200, right: 200 },
          children: [new TextRun({ text: "Help build Africa's interoperability layer. Join the community shaping how the continent exchanges data.", bold: true, color: "FFFFFF", size: 24 })] }),
        spacer(120),
        new Paragraph({ children: [new TextRun({ text: "This is a working draft (v0.1) published for community review. It will evolve in the open.", italics: true, color: MUTED, size: 19 })] }),
      ],
    },
  ],
});

Packer.toBuffer(doc).then((buf) => {
  fs.writeFileSync("IXEA-Whitepaper.docx", buf);
  console.log("wrote IXEA-Whitepaper.docx");
});
