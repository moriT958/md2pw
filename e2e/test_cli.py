"""E2E tests for md2pw CLI."""

import subprocess
import tempfile
from pathlib import Path

import pytest

# Path to project root
PROJECT_ROOT = Path(__file__).parent.parent
MD2PW_BINARY = PROJECT_ROOT / "md2pw"
SAMPLE_MD = PROJECT_ROOT / "testdata" / "sample.md"


@pytest.fixture(scope="session", autouse=True)
def build_binary():
    """Build the md2pw binary before running tests."""
    subprocess.run(
        ["go", "build", "-o", str(MD2PW_BINARY), "./cmd/cli"],
        cwd=PROJECT_ROOT,
        check=True,
    )
    yield
    # Cleanup after all tests
    if MD2PW_BINARY.exists():
        MD2PW_BINARY.unlink()


class TestFileInput:
    """Test: ./md2pw testdata/sample.md"""

    def test_converts_file_to_pukiwiki(self, build_binary):
        """sample.md を引数で渡して変換し、標準出力に PukiWiki 形式で出力される。"""
        result = subprocess.run(
            [str(MD2PW_BINARY), str(SAMPLE_MD)],
            capture_output=True,
            text=True,
            check=True,
        )
        output = result.stdout

        # H1 → * H1
        assert "* H1" in output

        # H2 → ** H2
        assert "** H2" in output

        # H3 → *** H3
        assert "*** H3" in output

        # H4 → そのまま (#### H4)
        assert "#### H4" in output

        # **Bold text** → ''Bold text''
        assert "''Bold text''" in output

        # - item → -item
        assert "-list1" in output
        assert "--list.a" in output

        # 1. item → +item
        assert "+ordered1" in output
        assert "++second" in output

        # コードブロック → 2スペースインデント
        assert '  fmt.Println("Hello")' in output

        # [text](url) → [[text>url]]
        assert "[[this is link>https://example.com]]" in output

        # テーブルヘッダー → |~ 形式
        assert "|~ Column1 |~ Column2 |~ Column3 |~ Column4 |" in output
        assert "| Item1.1 | Item2.1 | Item3.1 | Item4.1 |" in output


class TestPipeInput:
    """Test: cat testdata/sample.md | ./md2pw"""

    def test_converts_stdin_to_pukiwiki(self, build_binary):
        """stdin から入力を受け取って変換する。"""
        with open(SAMPLE_MD) as f:
            sample_content = f.read()

        result = subprocess.run(
            [str(MD2PW_BINARY)],
            input=sample_content,
            capture_output=True,
            text=True,
            check=True,
        )
        output = result.stdout

        # Same verifications as file input
        assert "* H1" in output
        assert "** H2" in output
        assert "*** H3" in output
        assert "#### H4" in output
        assert "''Bold text''" in output
        assert "-list1" in output
        assert "+ordered1" in output
        assert '  fmt.Println("Hello")' in output
        assert "[[this is link>https://example.com]]" in output
        assert "|~ Column1 |~ Column2 |~ Column3 |~ Column4 |" in output


class TestFileOutput:
    """Test: ./md2pw -o output.txt testdata/sample.md"""

    def test_writes_output_to_file(self, build_binary):
        """-o オプションでファイルに出力し、内容を検証する。"""
        with tempfile.TemporaryDirectory() as tmpdir:
            output_file = Path(tmpdir) / "output.txt"

            subprocess.run(
                [str(MD2PW_BINARY), "-o", str(output_file), str(SAMPLE_MD)],
                check=True,
            )

            assert output_file.exists()
            output = output_file.read_text()

            # Same verifications
            assert "* H1" in output
            assert "** H2" in output
            assert "*** H3" in output
            assert "#### H4" in output
            assert "''Bold text''" in output
            assert "-list1" in output
            assert "+ordered1" in output
            assert '  fmt.Println("Hello")' in output
            assert "[[this is link>https://example.com]]" in output
            assert "|~ Column1 |~ Column2 |~ Column3 |~ Column4 |" in output
