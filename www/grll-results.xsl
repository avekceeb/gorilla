<?xml version="1.0" encoding="utf-8" ?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform" >
<xsl:output omit-xml-declaration="yes" indent="yes"/>

<xsl:template match="/GrllHistorical">

    <xsl:variable name="verdict">
        <xsl:choose>
        <xsl:when test="count(results/status[text() = 'failed']) = 0">passed</xsl:when>
        <xsl:otherwise>failed</xsl:otherwise>
        </xsl:choose>
    </xsl:variable>


<html><head><title></title>
<style type="text/css">
    <!--/*--><![CDATA[/*><!--*/
    body { font-family: monospace ; }
    span {margin: 0px 10px ; }
    td { padding: 3px ; }
    tr.header { background-color: silver; }
    tr.odd { background-color: #E4E4E4; }
    tr.even { background-color: #F0F0F0 ; }
    .passed { color: green; }
    .failed { color: red; }
    .skipped { color: darkorange; }
    /*]]>*/-->
</style>
</head>
<body>
    <h1><xsl:value-of select="test"/></h1>
    <h4>
        <span class="passed">Passed = <xsl:value-of select="count(items/status[text() = 'passed'])"/></span>
        <span class="failed">Failed = <xsl:value-of select="count(items/status[text() = 'failed'])"/></span>
        <span class="skipped">Skipped = <xsl:value-of select="count(items/status[text() = 'skipped'])"/></span>

    </h4>

    <h3>Results</h3>
    <table cellspacing="0">
    <tr class="header">
        <td>Run</td>
        <td>Status</td>
        <td>Message</td>
        <td>Date</td>
    </tr>
    <xsl:apply-templates select='items' />
    </table>
</body>
</html>
</xsl:template>


<xsl:template match='items'>
    <xsl:variable name="class">
        <xsl:choose>
        <xsl:when test="normalize-space(./status) = 'passed'">passed</xsl:when>
        <xsl:when test="normalize-space(./status) = 'failed'">failed</xsl:when>
        <xsl:when test="normalize-space(./status) = 'skipped'">skipped</xsl:when>
        <xsl:otherwise>failed</xsl:otherwise>
        </xsl:choose>
    </xsl:variable>
    <xsl:variable name="oddeven">
        <xsl:choose>
        <xsl:when test="position() mod 2 = 0">even</xsl:when>
        <xsl:otherwise>odd</xsl:otherwise>
        </xsl:choose>
    </xsl:variable>
    <tr class="{$oddeven}">
        <td class="{$class}"><xsl:value-of select="normalize-space(./run)"/></td>
        <td class="{$class}"><xsl:value-of select="normalize-space(./status)"/></td>
        <td class="{$class}"><xsl:value-of select="normalize-space(./msg)"/></td>
        <td class="{$class}"><xsl:value-of select="normalize-space(./ts)"/></td>
    </tr>
</xsl:template>

</xsl:stylesheet>

